// Command agent is a CLI for the CodeAct-style file/dev assistant: it turns
// a natural-language task about a local directory into Go code, executes
// that code with an embedded interpreter, and retries with the execution
// feedback if it fails.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"codeact-agent/internal/agent"
	"codeact-agent/internal/executor"
	"codeact-agent/internal/llm"
	"codeact-agent/internal/tools"
)

func main() {
	var (
		workDir     = flag.String("workdir", ".", "directory the agent is allowed to operate on")
		maxAttempts = flag.Int("max-attempts", 3, "max generate/execute/fix attempts per task")
		verbose     = flag.Bool("v", false, "print generated code and raw output for every attempt")
	)
	flag.Parse()

	sandbox, err := tools.NewSandbox(*workDir)
	if err != nil {
		fatal("invalid -workdir: %v", err)
	}

	provider, err := llm.FromEnv()
	if err != nil {
		fatal("%v", err)
	}
	fmt.Fprintf(os.Stderr, "using model provider: %s\n", provider.Name())

	exec := executor.New()

	onStep := func(s agent.Step) {
		if !*verbose {
			return
		}
		fmt.Fprintf(os.Stderr, "\n--- attempt %d ---\n%s\n", s.Attempt, s.Code)
		if s.ExecErr != nil {
			fmt.Fprintf(os.Stderr, "[executor error] %v\n", s.ExecErr)
			return
		}
		if s.Output.Stdout != "" {
			fmt.Fprintf(os.Stderr, "[stdout] %s\n", s.Output.Stdout)
		}
		if s.Output.ToolError != "" {
			fmt.Fprintf(os.Stderr, "[tool error] %s\n", s.Output.ToolError)
		}
	}

	ag := agent.New(provider, exec, sandbox, *maxAttempts, onStep)

	tasks := flag.Args()
	if len(tasks) > 0 {
		if !runTask(ag, strings.Join(tasks, " ")) {
			os.Exit(1)
		}
		return
	}

	repl(ag)
}

// runTask executes a single task and prints its result. It returns whether
// the task succeeded; callers decide what to do with a failure (exit code
// for one-shot mode, keep prompting for the REPL).
func runTask(ag *agent.Agent, task string) bool {
	outcome, err := ag.Do(context.Background(), task)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return false
	}
	if !outcome.Success {
		last := outcome.Steps[len(outcome.Steps)-1]
		if last.ExecErr != nil {
			fmt.Fprintf(os.Stderr, "failed after %d attempts: %v\n", len(outcome.Steps), last.ExecErr)
		} else {
			fmt.Fprintf(os.Stderr, "failed after %d attempts: %s\n", len(outcome.Steps), last.Output.ToolError)
		}
		return false
	}
	fmt.Println(outcome.FinalAnswer)
	return true
}

func repl(ag *agent.Agent) {
	fmt.Println("CodeAct file/dev assistant. Type a task, or 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			return
		}
		runTask(ag, line)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
