// Package agent implements the generate -> execute -> observe -> fix loop
// that turns a natural-language task into executed Go code.
package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"codeact-agent/internal/executor"
	"codeact-agent/internal/llm"
	"codeact-agent/internal/tools"
)

// Step records one iteration of the loop, for logging/inspection.
type Step struct {
	Attempt int
	Code    string
	Output  executor.Result
	ExecErr error // executor-level failure (compile error, panic, timeout)
}

// Outcome is the final result of running a task to completion (success or
// exhausted retries).
type Outcome struct {
	Steps   []Step
	Success bool
	// FinalAnswer is the string returned by the last successful Run().
	FinalAnswer string
}

// Agent wires a model provider to the sandboxed executor.
type Agent struct {
	provider    llm.Provider
	executor    *executor.Executor
	sandbox     *tools.Sandbox
	maxAttempts int
	onStep      func(Step)
}

// New builds an Agent scoped to sandbox for the lifetime of every Do() call.
// onStep, if non-nil, is called after every attempt so callers (e.g. the
// CLI) can stream progress; it may be nil.
func New(provider llm.Provider, exec *executor.Executor, sandbox *tools.Sandbox, maxAttempts int, onStep func(Step)) *Agent {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Agent{provider: provider, executor: exec, sandbox: sandbox, maxAttempts: maxAttempts, onStep: onStep}
}

// Do runs the generate/execute/fix loop for a single task description and
// returns once the code succeeds or attempts are exhausted.
func (a *Agent) Do(ctx context.Context, task string) (Outcome, error) {
	var (
		steps       []Step
		lastCode    string
		failureNote string
		haveFailure bool
	)

	for attempt := 1; attempt <= a.maxAttempts; attempt++ {
		var prompt string
		if !haveFailure {
			prompt = userPrompt(task)
		} else {
			prompt = retryPrompt(task, lastCode, failureNote)
		}

		raw, err := a.provider.Generate(ctx, systemPrompt, prompt)
		if err != nil {
			return Outcome{Steps: steps}, fmt.Errorf("model generation failed: %w", err)
		}

		code, err := extractCode(raw)
		if err != nil {
			step := Step{Attempt: attempt, ExecErr: err}
			steps = append(steps, step)
			a.report(step)
			lastCode = raw
			failureNote = err.Error()
			haveFailure = true
			continue
		}
		lastCode = code

		result, execErr := a.executor.Run(ctx, a.sandbox, code)
		step := Step{Attempt: attempt, Code: code, Output: result, ExecErr: execErr}
		steps = append(steps, step)
		a.report(step)

		if execErr != nil {
			failureNote = execErr.Error()
			haveFailure = true
			continue
		}
		if result.ToolError != "" {
			failureNote = "Run() returned an error: " + result.ToolError
			haveFailure = true
			continue
		}

		return Outcome{Steps: steps, Success: true, FinalAnswer: result.ReturnValue}, nil
	}

	return Outcome{Steps: steps, Success: false}, nil
}

func (a *Agent) report(s Step) {
	if a.onStep != nil {
		a.onStep(s)
	}
}

var codeBlockRE = regexp.MustCompile("(?s)```(?:go)?\\s*\\n(.*?)```")

// extractCode pulls the first fenced code block out of a model response. If
// no fence is present but the response looks like bare Go source (starts
// with "package "), it is used as-is.
func extractCode(raw string) (string, error) {
	if m := codeBlockRE.FindStringSubmatch(raw); m != nil {
		return strings.TrimSpace(m[1]), nil
	}
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "package ") {
		return trimmed, nil
	}
	return "", fmt.Errorf("no Go code block found in model response: %q", truncate(raw, 200))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
