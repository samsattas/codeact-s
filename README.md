# CodeAct Agent

A Go agent that solves tasks about a local project directory — listing
files, searching for patterns, summarizing a codebase, editing files,
running `go build`/`go test`/`git` — by writing and executing real Go code
against a sandboxed API, instead of returning a fixed JSON tool call. Ask it
something in plain English and watch it write, run, and self-correct actual
Go source in real time.

**Live demo:** _pending deployment — see [Deployment](#deployment) below._

## What it does

- Understands natural-language tasks about a project directory: "find all
  TODO comments and summarize them," "count lines of code by file
  extension," "run `go vet` and tell me if there's anything to fix."
- Generates a complete Go function to accomplish the task, executes it
  immediately in an embedded interpreter, and reads back the result.
- If the code fails to compile or a step returns an error, the exact failure
  is fed back to the model automatically, and it gets another attempt — up
  to a configurable limit — before giving up.
- Ships with two interfaces to the same agent: a terminal CLI and a browser
  UI that streams every attempt live (generated code, output, pass/fail)
  as it happens.
- Works with a local model via [Ollama](https://ollama.com) out of the box,
  or a cloud model via the Anthropic API if you'd rather.

## How it works

Every task is answered by generating a single Go function —
`func Run() (string, error)` — that calls into a small, fixed set of
sandboxed operations (read/write files, list a directory, grep, walk a file
tree, run an allowlisted command). That function is compiled and executed
on the spot by an embedded Go interpreter, not sent to a shell or written to
disk as a script. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the
full design: the sandboxing model, why an embedded interpreter instead of a
subprocess, and the retry loop that lets the agent recover from its own
mistakes.

## Getting started

Requires Go 1.23+.

```bash
go build -o bin/agent ./cmd/agent
go build -o bin/web ./cmd/web
```

### Pick a model

**Local, no API key (default)** — install [Ollama](https://ollama.com), pull
a coding-capable model, and run it:

```bash
ollama serve
ollama pull qwen2.5-coder:7b
```

**Cloud (Anthropic)** — set an API key and it takes priority automatically:

```bash
export ANTHROPIC_API_KEY=sk-...
```

### Run the CLI

```bash
./bin/agent -workdir ./some/project "list the .go files and count them"
```

Flags: `-workdir` (sandbox root, default `.`), `-max-attempts` (default 3),
`-v` (print each attempt's generated code and result). Run with no task
argument to get an interactive REPL.

### Run the web UI

```bash
./bin/web -workdir ./some/project -addr :8080
```

Open `http://localhost:8080`, type a task (or click an example), and hit
Run. It's a single self-contained Go binary — the page is embedded in the
binary at build time, no separate asset server or JS build step.

`-workdir` only sets the *default* directory. The page has its own
**Dir** field, pre-filled with that default but editable per run, so a
single running server can point different runs at different directories —
it's no longer tied to whichever one it started with. There's no
authentication and no allowlist of directories, so treat the web UI as a
single-user local tool, not something to expose beyond `localhost`.

## Testing

```bash
go test ./...
```

The full suite runs offline — no network access or running model required —
using a scripted fake model provider to exercise the generate/execute/fix
loop deterministically. See [docs/TESTING.md](docs/TESTING.md) for what's
covered and current coverage numbers.

## Tech stack

Go · [yaegi](https://github.com/traefik/yaegi) (embedded Go interpreter) ·
Ollama / Anthropic API · vanilla HTML/CSS/JS (no build step, no framework)

## Project layout

```
cmd/agent/        CLI entrypoint
cmd/web/          Browser UI, served over HTTP
internal/agent/   The generate -> execute -> observe -> fix loop
internal/executor/  Runs generated Go code via an embedded interpreter
internal/tools/   The sandboxed API generated code is allowed to call
internal/llm/     Model provider interface (Ollama, Anthropic)
docs/             Architecture and testing notes
```

## Deployment

Not yet deployed. This section will link to a live instance once one is up.
