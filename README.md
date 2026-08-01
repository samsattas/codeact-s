# CodeAct File/Dev Agent

A small proof-of-concept agent, written in Go, that follows the
**CodeAct ("code as action")** paradigm: instead of asking a model to emit a
JSON tool call, it asks the model to write a real Go function, and then
**executes that Go code** with an embedded interpreter
([yaegi](https://github.com/traefik/yaegi)) to perform the action.

The agent's job: answer questions and perform tasks about a local project
directory — list files, grep for patterns, summarize a codebase, edit files,
run `go build`/`go test`/`git`, etc. — by writing and running Go snippets
against a small sandboxed tools API.

## Why this counts as "code as action"

A classic tool-calling agent gets a JSON blob like
`{"tool": "grep", "pattern": "TODO"}` back from the model, and the harness
interprets it. This agent instead gets back something like:

```go
package main

import (
	"fmt"
	"tools"
)

func Run() (string, error) {
	matches, err := tools.Grep(".", "TODO")
	if err != nil {
		return "", err
	}
	summary := fmt.Sprintf("found %d TODOs", len(matches))
	for _, m := range matches {
		summary += fmt.Sprintf("\n%s:%d: %s", m.File, m.Line, m.Text)
	}
	return summary, nil
}
```

...and **runs it as actual Go**, with loops, conditionals, error handling,
and multiple tool calls composed in a single action — not a single flat
function call. If it fails (compile error, or the tool call itself returns
an error), the failure is fed back to the model as an execution trace and it
gets another attempt to fix its own code. That generate → execute → observe
→ fix loop is the core of CodeAct, and it's implemented literally here: the
"action" the agent takes on the environment *is* compiled and run Go source,
not a call dispatched from a fixed schema.

## Architecture

```
cmd/agent/main.go        CLI entrypoint (flags, one-shot or REPL mode)
internal/llm/             Provider interface + Ollama and Anthropic backends
internal/agent/           The generate -> execute -> observe -> fix loop
internal/executor/        Runs generated Go code via yaegi (embedded interpreter)
internal/tools/           The sandboxed API generated code is allowed to call
```

- **`internal/tools`** implements file read/write/list, a directory tree
  renderer, grep, a line-counter, and a whitelisted command runner
  (`go`, `git`, `ls`, `wc`, `find`, `cat`, `echo`, `gofmt`). Every path is
  resolved against a single sandbox root set via `tools.SetWorkDir` and
  rejected if it would escape that root (`..`, absolute paths, etc.).
- **`internal/executor`** spins up a fresh `yaegi` interpreter per attempt,
  loads only a curated subset of the Go standard library
  (`fmt`, `strings`, `strconv`, `errors`, `sort`, `time`, `regexp`, `bytes`,
  `unicode`, `math`, `encoding/json`) plus the `tools` package described
  above, evaluates the generated source, and calls `Run()`. Deliberately
  **excluded**: `os`, `io`, `io/ioutil`, `os/exec`, `net/http`, `syscall`.
  This was not a hypothetical concern — see [Sandbox design](#sandbox-design)
  below for the real bug this caught during testing.
- **`internal/agent`** owns the retry loop: build a prompt, call the model,
  extract the fenced Go code block from the response, run it, and if it
  fails, build a follow-up prompt that includes the previous code and the
  exact failure so the model can fix it (up to `-max-attempts`, default 3).
- **`internal/llm`** is a one-method `Provider` interface
  (`Generate(ctx, systemPrompt, userPrompt) (string, error)`) so swapping
  models never touches the agent or executor. `llm.FromEnv()` picks
  Anthropic if `ANTHROPIC_API_KEY` is set, otherwise falls back to a local
  Ollama server.

## Sandbox design

While testing against a real local model, the first working version of the
executor loaded yaegi's **entire** standard library, including `os` and
`io/ioutil`. The model happily generated code that called `ioutil.ReadDir(".")`
directly instead of `tools.ListDir(".")` — which compiled and ran fine, but
`"."` inside the interpreter resolves to the *agent process's* real working
directory, completely bypassing the `-workdir` sandbox. That's a real
example of why "the LLM is just a tool, you own the correctness" matters:
the fix was to stop trusting the model to only use the safe API by
*convention*, and instead make the unsafe packages structurally unavailable
by only registering an allowlisted subset of the stdlib with the
interpreter (`internal/executor/executor.go`, `allowedStdlibPackages`).
Now `import "os"` fails to compile inside the sandbox, full stop — there is
no filesystem/process/network access generated code can reach that doesn't
go through `internal/tools`, which itself enforces the sandbox root.

Other safety measures:
- `tools.RunCommand` only allows an explicit binary allowlist, never goes
  through a shell (`exec.Command(name, args...)`, no `sh -c`), so argument
  values can't inject shell metacharacters, and it has a 20s hard timeout.
- Every generated snippet runs against a **fresh interpreter instance** —
  no state leaks between attempts or between tasks.
- Each execution has a soft 10s timeout. Known limitation: yaegi can't
  forcibly preempt a runaway goroutine (e.g. `for {}`), so a truly
  pathological snippet leaks a goroutine until the process exits rather
  than being killed outright. A production version of this would run
  generated code in a subprocess (or a container) instead of in-process.

## Build & run

Requires Go 1.23+.

```bash
go build -o bin/agent ./cmd/agent
```

### Model backend

Pick one:

**Local (Ollama, default, no API key needed)**
```bash
ollama serve                        # in another terminal
ollama pull qwen2.5-coder:7b        # or any other coder-capable model
./bin/agent -workdir ./some/project "list the .go files and count them"
```
Override the model/host with `OLLAMA_MODEL` / `OLLAMA_HOST` env vars.

**Cloud (Anthropic)**
```bash
export ANTHROPIC_API_KEY=sk-...
./bin/agent -workdir ./some/project "summarize the TODOs in this project"
```
Override the model with `ANTHROPIC_MODEL` (default `claude-sonnet-4-5`).
If `ANTHROPIC_API_KEY` is set it always takes priority over Ollama.

### CLI flags

```
-workdir string       sandbox root the agent may read/write (default ".")
-max-attempts int     generate/execute/fix attempts per task (default 3)
-v                    print generated code + execution trace for every attempt
```

Run with no task argument for an interactive REPL:
```bash
./bin/agent -workdir ./some/project -v
> find all TODO comments and write a summary to todos.md
> exit
```

### Tests

```bash
go test ./...
```

Covers the tools sandbox (path escape, allowlisted commands, grep/count
correctness), the executor (successful runs, compile errors, tool errors,
timeouts, the `os`-is-unavailable guarantee), and the agent retry loop using
a scripted fake `llm.Provider` — no network or real model required to run
`go test ./...`.

## Real transcript (local model, `qwen2.5-coder:3b` via Ollama)

This is an unedited run against a tiny local model, included because it
shows both halves of the CodeAct loop working: a genuine compile failure
being self-corrected, and — just as importantly — a case where the code ran
successfully but the model's *logic* was still slightly wrong (it wrote a
matching `-l` per file it iterated over instead of using `match.File`),
producing a duplicate line in the output. The harness catches and retries
technical failures (compile errors, failed tool calls); it does **not**
validate semantic correctness of a successful run — that's an inherent
limitation of the approach, not a bug, and worth being upfront about.

```
$ ./bin/agent -workdir ./demo -v "list the files in this directory and tell me how many .go files there are"
using model provider: ollama:qwen2.5-coder:3b

--- attempt 1 ---
... (imports "tools" and "fmt" but uses strings.ToLower without importing "strings")
[executor error] compile error: 16:10: undefined: strings

--- attempt 2 ---
... (same code, now with "strings" imported)
Directory contents:
demo/
├── main.go
└── math.go

Number of .go files: 0
```

## Design choices (short version)

- **yaegi over a subprocess/shell script**: keeps "code as action" literally
  true to the Go challenge — the model writes and runs *Go*, not a bash
  script wrapping the language you were asked to build in. It also means no
  external interpreter binary is required at runtime, and lets the sandbox
  be enforced at the symbol-table level (see above) instead of trying to
  sandbox an OS process.
- **A tiny, fixed tools API instead of general stdlib access**: the whole
  point of CodeAct is composing *known-safe* operations with real control
  flow, not giving a model a shell. Fewer, well-typed tools also make
  self-correction easier for small/local models, which struggle more with
  large ambiguous surfaces.
- **Ollama-by-default**: the challenge can be demoed/re-run offline, without
  provisioning or exposing an API key, while `ANTHROPIC_API_KEY` is a
  one-line opt-in to a much stronger model.

## What's intentionally out of scope (PoC, not production)

- No persistent conversation memory across separate CLI invocations.
- No true kill-on-timeout for runaway generated code (see above).
- No multi-file / multi-package generated programs — one `Run()` per task.
- No streaming output from the model while it's generating.
