# Testing

```bash
go test ./...            # run everything
go test ./... -cover     # with coverage
go test ./... -v         # verbose, per-test output
```

The full suite runs offline: no network access and no running model
required. The agent's generate/execute/fix loop is exercised with a
scripted fake model provider (a `Provider` implementation that returns a
fixed sequence of canned responses per test), so the retry logic is tested
deterministically instead of depending on what a real model happens to
output on a given run.

## Coverage by package

| Package | Coverage | What's exercised |
|---|---|---|
| `internal/tools` | 82.1% | Path-escape rejection, read/write/list, grep correctness, line counting (including that compiled binaries in `bin/` aren't miscounted as text), the `RunCommand` binary allowlist |
| `internal/agent` | 76.0% | First-try success, retry-and-recover after a compile error, retry-and-recover after a tool error, exhausting all attempts, code-block extraction from a model response |
| `internal/executor` | 75.6% | Successful execution, compile errors, tool errors surfaced through `Run()`, that `import "os"` is unavailable inside the sandbox, that paths outside the sandbox root are rejected |
| `cmd/web` | 56.9% | `/api/info`, `/api/run` end to end (streamed NDJSON success and failure), rejecting non-`POST` and malformed requests |

`cmd/agent` and `internal/llm` are not covered by automated tests: the
former is a thin CLI wrapper with no branching logic of its own, and the
latter makes real HTTP calls to Ollama/Anthropic — its correctness is
validated by exercising it directly (see below) rather than by mocking the
HTTP layer.

## What the security-relevant tests actually check

Two tests exist specifically because a real issue was found while building
this, not as a hypothetical:

- `TestStdlibOSPackageIsUnavailable` (`internal/executor`) — asserts that
  generated code importing `os` fails to compile inside the interpreter.
  An earlier version exposed the full standard library, which let generated
  code bypass the sandbox entirely by calling `ioutil.ReadDir(".")`
  directly instead of going through `tools.ListDir`.
- `TestSandboxBlocksPathEscape` / `TestResolveBlocksEscape` — assert that a
  path like `../../etc/passwd` is rejected rather than resolved outside the
  sandbox root.

## Manual end-to-end validation

Beyond the automated suite, the CLI and web UI were both run against real
local models (`qwen2.5-coder:1.5b/3b/7b` via Ollama) on real tasks — listing
files, counting lines of code by extension, grepping for TODO comments and
writing a summary file, running `go vet` — to confirm the full loop behaves
correctly against non-scripted model output, including the self-correction
path after a genuine compile failure. See
[docs/ARCHITECTURE.md](ARCHITECTURE.md#observed-behavior) for an example
transcript and the limitation that follows from it.
