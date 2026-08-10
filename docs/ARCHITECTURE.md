# Architecture

## Why this counts as "code as action"

A classic tool-calling agent gets a JSON blob like
`{"tool": "grep", "pattern": "TODO"}` back from the model, and a harness
interprets it. This agent instead gets back a complete Go function:

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

...and runs it as actual Go, with loops, conditionals, error handling, and
multiple tool calls composed in a single action — not a single flat
function call. If it fails (compile error, or the tool call itself returns
an error), the failure is fed back to the model as an execution trace and it
gets another attempt to fix its own code. That generate → execute → observe
→ fix loop is the core of CodeAct: the action the agent takes on the
environment *is* compiled and run Go source, not a call dispatched from a
fixed schema.

## Module layout

```
cmd/agent/main.go   CLI entrypoint (flags, one-shot or REPL mode)
cmd/web/main.go     Web UI: same agent, served over HTTP with a streamed browser view
internal/llm/       Provider interface + Ollama and Anthropic backends
internal/agent/     The generate -> execute -> observe -> fix loop
internal/executor/  Runs generated Go code via yaegi (embedded interpreter)
internal/tools/     The sandboxed API generated code is allowed to call
```

`cmd/agent` and `cmd/web` are both thin front-ends over the same
`internal/agent` + `internal/executor` + `internal/tools` core — neither one
duplicates the loop; they call `agent.New(...).Do(...)` and render the
result differently (stdout vs. a browser page).

- **`internal/tools`** implements file read/write/list, a directory tree
  renderer, grep, a line-counter, and an allowlisted command runner (`go`,
  `git`, `ls`, `wc`, `find`, `cat`, `echo`, `gofmt`). Every path is resolved
  against a single sandbox root set via `tools.SetWorkDir` and rejected if
  it would escape that root (`..`, absolute paths, etc).
- **`internal/executor`** spins up a fresh `yaegi` interpreter per attempt,
  loads only a curated subset of the Go standard library (`fmt`, `strings`,
  `strconv`, `errors`, `sort`, `time`, `regexp`, `bytes`, `unicode`, `math`,
  `encoding/json`) plus the `tools` package above, evaluates the generated
  source, and calls `Run()`. Deliberately excluded: `os`, `io`,
  `io/ioutil`, `os/exec`, `net/http`, `syscall` — see
  [Sandbox design](#sandbox-design).
- **`internal/agent`** owns the retry loop: build a prompt, call the model,
  extract the fenced Go code block from the response, run it, and on
  failure build a follow-up prompt with the previous code and the exact
  error so the model can fix it, up to `-max-attempts` (default 3).
- **`internal/llm`** is a one-method `Provider` interface
  (`Generate(ctx, systemPrompt, userPrompt) (string, error)`), so swapping
  models never touches the agent or executor. `llm.FromEnv()` selects
  Anthropic if `ANTHROPIC_API_KEY` is set, otherwise a local Ollama server.

## Sandbox design

Generated code only has access to the sandboxed API in `internal/tools` —
not the real filesystem, process list, or network. That guarantee is
enforced at two levels:

1. **Symbol availability.** `internal/executor` registers only an
   allowlisted subset of the Go standard library with the interpreter.
   Packages that reach outside the sandbox — `os`, `io`, `io/ioutil`,
   `os/exec`, `net/http`, `syscall` — are structurally unavailable:
   `import "os"` fails to compile inside the interpreter, full stop. There
   is no filesystem, process, or network access generated code can reach
   that doesn't go through `internal/tools`.
2. **Path confinement.** Every function in `internal/tools` resolves its
   path argument against a single sandbox root (set once via
   `tools.SetWorkDir`) and rejects anything that would resolve outside it
   (`..`, absolute paths, symlink escapes via `filepath.Clean` + `Rel`).

Other safety measures:
- `tools.RunCommand` only allows an explicit binary allowlist, and never
  goes through a shell (`exec.Command(name, args...)`, no `sh -c`), so
  argument values can't inject shell metacharacters. It has a 20s hard
  timeout.
- Every generated snippet runs against a fresh interpreter instance — no
  state leaks between attempts or between tasks.
- Each execution has a soft 10s timeout via `context` + `select`. Known
  limitation: the interpreter can't forcibly preempt a runaway goroutine
  (e.g. a generated `for {}`), so a pathological snippet leaks a goroutine
  until the process exits rather than being killed outright. A production
  version of this would run generated code in a subprocess or container
  instead of in-process, where the OS can kill it.

## Design choices

- **An embedded Go interpreter ([yaegi](https://github.com/traefik/yaegi))
  over a subprocess or shell script.** The model writes and runs real Go,
  not a script wrapping it. It also means no external interpreter binary is
  required at runtime, and the sandbox can be enforced at the symbol-table
  level instead of trying to contain an OS process.
- **A tiny, fixed tools API instead of general standard-library access.**
  The point of code-as-action is composing known-safe operations with real
  control flow, not exposing a shell. A smaller, well-typed surface also
  makes self-correction more tractable for smaller/local models.
- **Ollama by default, Anthropic as an opt-in.** The agent runs fully
  offline with no API key required; setting `ANTHROPIC_API_KEY` switches to
  a stronger cloud model with no code changes.

## Observed behavior

Smaller local models occasionally succeed on the first attempt and
occasionally need one or two retries — e.g. generating code that calls
`strings.HasSuffix` without importing `"strings"`, then correcting itself
once the compile error is fed back:

```
--- attempt 1 ---
... (imports "tools" and "fmt" but uses strings.ToLower without importing "strings")
[executor error] compile error: 16:10: undefined: strings

--- attempt 2 ---
... (same code, now with "strings" imported)
Directory contents:
demo/
├── main.go
└── math.go
```

The retry loop catches and resolves *technical* failures — compile errors,
failed tool calls — using the interpreter's own error output as ground
truth. It does not, and cannot, validate the *semantic* correctness of a
successful run (e.g. code that compiles and returns an answer, but counted
the wrong thing). That's an inherent limitation of the approach worth
stating plainly rather than glossing over.

## Known limitations (proof of concept, not production)

- No persistent conversation memory across separate CLI invocations.
- No true kill-on-timeout for runaway generated code (see Sandbox design).
- No multi-file / multi-package generated programs — one `Run()` per task.
- No token-level streaming from the model itself; both front-ends show a
  complete attempt (code + result) once it finishes, not word-by-word.
- The web UI has no authentication and binds to a single fixed sandbox root
  for the lifetime of the process.
