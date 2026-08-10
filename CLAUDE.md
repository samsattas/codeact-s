# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).

## Commands

```bash
go build -o bin/agent ./cmd/agent      # build the CLI
go build -o bin/web ./cmd/web          # build the web UI (page is embedded via go:embed)

go test ./...                          # run the full suite (offline, no model/network needed)
go test ./... -cover                   # with coverage
go test ./internal/agent/ -run TestAgentRetriesAfterCompileErrorThenSucceeds -v   # a single test

go vet ./...
```

Run the CLI: `./bin/agent -workdir ./some/project "task"` (one-shot) or with no
task argument for a REPL. Flags: `-workdir` (sandbox root, default `.`),
`-max-attempts` (default 3), `-v` (print each attempt's code/output).

Run the web UI: `./bin/web -workdir ./some/project -addr :8080`, then open
`http://localhost:8080`. `/api/run` streams NDJSON (one `{"type":"step",...}`
line per attempt, then a final `{"type":"done",...}`); `/api/info` reports the
active provider, sandbox root, and max attempts.

Model backend: Ollama by default (`ollama serve` + `ollama pull
qwen2.5-coder:7b`, no key needed) — set `ANTHROPIC_API_KEY` to switch to
Claude instead (`internal/llm.FromEnv` prefers it automatically when set).

## Architecture

This is a CodeAct-style agent: instead of returning a fixed JSON tool call,
the model returns one complete Go function —
`func Run() (string, error)` — which is compiled and executed on the spot by
an embedded interpreter ([yaegi](https://github.com/traefik/yaegi)), not sent
to a shell or written to disk. If it fails to compile or a step returns an
error, the exact failure is fed back to the model as the next prompt (see
`internal/agent/prompt.go:retryPrompt` and `diagnosisHint`, which pattern-matches
common failure shapes — undefined symbols, unused imports, a RunCommand
non-zero exit — into a specific corrective instruction rather than just
echoing the raw compiler error) and it gets another attempt, up to
`-max-attempts`.

`cmd/agent` (CLI) and `cmd/web` (browser UI, static assets embedded via
`//go:embed static`) are both thin front-ends over the same core loop —
neither duplicates it, they both just call `agent.New(...).Do(...)` and
render the result differently (stdout vs. a streamed NDJSON page):

- **`internal/agent`** owns the generate → execute → observe → fix loop
  (`Agent.Do`): build a prompt, call the model, extract the fenced Go code
  block, run it via the executor, and on failure build a retry prompt with
  the previous code + the exact error (plus a targeted hint) so the model can
  fix its own mistake.
- **`internal/executor`** spins up a fresh `yaegi` interpreter per attempt,
  loads only a curated stdlib subset (fmt, strings, strconv, errors, sort,
  time, regexp, bytes, unicode, math, encoding/json) plus the `tools`
  package, evaluates the source, and calls `Run()`.
- **`internal/tools`** is the sandboxed API generated code is allowed to
  call: read/write/list files, a directory tree renderer, grep, a
  line-counter, and an allowlisted command runner (`go`, `git`, `ls`, `wc`,
  `find`, `cat`, `echo`, `gofmt`, never through a shell). Every path is
  resolved against a single sandbox root (`tools.SetWorkDir`) and rejected if
  it would escape that root.
- **`internal/llm`** is a one-method `Provider` interface
  (`Generate(ctx, systemPrompt, userPrompt) (string, error)`); `FromEnv()`
  picks Anthropic or Ollama.

**Sandbox enforcement is two-layered, not just prompt-level convention**:
`internal/executor` only registers an allowlisted stdlib subset with the
interpreter, so `import "os"`/`io`/`os/exec`/`net/http`/`syscall` fails to
compile inside the interpreter regardless of what the model tries; separately,
every `internal/tools` function resolves its path argument against the
sandbox root and rejects anything that would resolve outside it. See
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full design rationale
(why an embedded interpreter over a subprocess, known limitations like no
true kill-on-timeout for runaway generated code) and
[docs/TESTING.md](docs/TESTING.md) for coverage details and which tests
specifically guard the sandbox boundary.
