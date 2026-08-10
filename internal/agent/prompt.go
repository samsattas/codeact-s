package agent

import (
	"fmt"
	"regexp"
)

// systemPrompt describes the tools API and the contract generated code must
// follow. It is sent once as the system message on every call.
const systemPrompt = `You are a coding agent that solves tasks about a local project directory by
writing and running real Go code (the "code as action" / CodeAct approach),
instead of calling tools through a fixed JSON schema.

You must respond with EXACTLY ONE Go code block (` + "```go ... ```" + `) and nothing else
outside of it. The code must:

  1. Declare "package main".
  2. Import "tools" (and any needed standard library packages) as needed.
  3. Define exactly one function: func Run() (string, error)
  4. Do all of the work for the task inside Run(), calling the tools API
     below. Use loops, conditionals, string building, error handling, etc.
     freely, exactly as you would in a normal Go program.
  5. Return a human-readable summary string as the first value. Return a
     non-nil error only if the task genuinely could not be completed.

Available tools package API (import "tools"):

  func tools.ReadFile(path string) (string, error)
      Reads a text file relative to the sandbox root.

  func tools.WriteFile(path string, content string) error
      Writes content to a file relative to the sandbox root, creating
      parent directories as needed.

  func tools.ListDir(path string) ([]string, error)
      Lists entries directly inside path (directories have a trailing "/").

  func tools.FileTree(path string, maxDepth int) (string, error)
      Returns a rendered directory tree, maxDepth levels deep.

  func tools.Grep(root string, pattern string) ([]tools.GrepMatch, error)
      Searches root recursively for a regular expression. GrepMatch has
      fields File string, Line int, Text string.

  func tools.CountLinesByExt(root string) (map[string]int, error)
      Returns a line count per file extension under root.

  func tools.RunCommand(name string, args ...string) (string, error)
      Runs an allowlisted command (go, git, ls, wc, find, cat, echo, gofmt)
      with the given args (no shell involved) and returns combined
      stdout+stderr as the first value. Its error return is non-nil
      whenever the command exits non-zero — routine for tools like
      "go vet", "go test", or "gofmt -l" when they find something to
      report, NOT a sign that something went wrong.

IMPORTANT — RunCommand's error is not your error. Run()'s own error means
"the task could not be completed"; it is a completely different thing from
RunCommand exiting non-zero. When a command's non-zero exit is itself the
expected result (a linter/vet/test finding issues, "git diff --exit-code"
finding changes, etc.), IGNORE RunCommand's error and return its output as
your answer with a NIL error from Run(). For example:

    output, _ := tools.RunCommand("go", "vet", "./...")
    if strings.TrimSpace(output) == "" {
        return "go vet found no issues.", nil
    }
    return "go vet found issues:\n" + output, nil

Note the "_" — RunCommand's error is deliberately discarded there, because
the output (not the exit code) is the actual answer to the task. Only
return a non-nil error from Run() if the command's output is empty or
clearly shows the command itself could not run (e.g. "not found", "no such
file").

Only these standard library packages are available: fmt, strings, strconv,
errors, sort, time, regexp, bytes, unicode, unicode/utf8, math,
encoding/json. Packages like os, io, io/ioutil, os/exec, and net/http are
NOT available and will fail to compile — use the tools.* functions above
for any filesystem, process, or network access instead.

Rules:
  - All paths are relative to the sandbox root ("." is the root itself).
  - Never try to access paths outside the sandbox root (e.g. "..", "/etc");
    the sandbox will reject them at runtime.
  - Do not invent tools.* functions beyond the ones listed above.
  - Do not import os, io, io/ioutil, os/exec, or net/http; they are not
    available. Use tools.* instead.
  - Keep the code focused on the task; do not add unrelated work.
  - fmt.Println inside Run() is fine for progress notes; the final answer
    should still be returned as the string result.
`

func userPrompt(task string) string {
	return fmt.Sprintf("Task: %s\n\nWrite the Go code now.", task)
}

func retryPrompt(task, previousCode, failureReason string) string {
	return fmt.Sprintf(
		`Task: %s

Your previous code failed. Read the failure message below carefully — it
tells you exactly what is wrong. Fix that specific problem and return a
corrected, complete code block following the exact same contract (package
main, import "tools", func Run() (string, error)). Do not return the same
code again; it must change to address the failure.

Previous code:
%s

Failure:
%s
%s
Write the corrected Go code now.`, task, fenced(previousCode), failureReason, diagnosisHint(failureReason))
}

func fenced(code string) string {
	return "```go\n" + code + "\n```"
}

var (
	undefinedRE    = regexp.MustCompile(`undefined: (\w+)`)
	unusedImportRE = regexp.MustCompile(`imported and not used`)
	unusedVarRE    = regexp.MustCompile(`declared (and|but) not used`)
	noNewVarsRE    = regexp.MustCompile(`no new variables on left side of :=`)
	exitStatusRE   = regexp.MustCompile(`exit status \d+`)
)

// diagnosisHint recognizes a handful of common yaegi/Go compiler error
// shapes and turns them into a blunt, specific instruction. Small models
// tend to re-generate the same broken code when only shown a bare compiler
// error; naming the exact fix dramatically improves retry success.
func diagnosisHint(failureReason string) string {
	if m := undefinedRE.FindStringSubmatch(failureReason); m != nil {
		name := m[1]
		return fmt.Sprintf(
			"\nHint: %q is undefined. This almost always means one of:\n"+
				"  - you used it without importing its package (e.g. used fmt.* but\n"+
				"    didn't `import \"fmt\"`) — check your import block lists every\n"+
				"    package the code references;\n"+
				"  - you misspelled it, or invented a tools.* function that doesn't\n"+
				"    exist — only use the exact functions and packages listed above.\n",
			name,
		)
	}
	if unusedImportRE.MatchString(failureReason) {
		return "\nHint: you imported a package but never referenced it in the code. Remove the unused import.\n"
	}
	if unusedVarRE.MatchString(failureReason) {
		return "\nHint: you declared a variable but never used it. Remove it, use it, or assign it to `_`.\n"
	}
	if noNewVarsRE.MatchString(failureReason) {
		return "\nHint: every variable on the left of `:=` already exists in this scope. Use `=` instead, or introduce at least one new variable name.\n"
	}
	if exitStatusRE.MatchString(failureReason) {
		return "\nHint: an \"exit status N\" error means tools.RunCommand's underlying command\n" +
			"exited non-zero — routine for `go vet`/`go test`/etc. reporting real findings,\n" +
			"not a bug in your code. RunCommand still gives you the command's combined\n" +
			"output alongside that error. If that output is itself the answer to the task,\n" +
			"return `output, nil` from Run() — do NOT do `return output, err` or\n" +
			"`return \"\", err`. A non-nil error from Run() means the task failed; forwarding\n" +
			"RunCommand's error there causes this exact retry loop, since the caller will\n" +
			"just keep asking you to fix a \"failure\" that was actually a valid result.\n"
	}
	return ""
}
