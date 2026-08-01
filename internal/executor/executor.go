// Package executor runs LLM-generated Go source code using yaegi, an
// embedded Go interpreter. This is the heart of the "code as action"
// mechanism: instead of parsing a JSON tool call, we execute real Go code
// that calls into the tools package directly.
package executor

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	"codeact-agent/internal/tools"
)

// Result carries everything the agent needs to decide whether an action
// succeeded and, if not, what to tell the model to fix.
type Result struct {
	// Stdout is anything the generated code printed (fmt.Println, etc.)
	// while it ran.
	Stdout string
	// ReturnValue is the string returned by the generated Run() function.
	ReturnValue string
	// ToolError is non-empty when Run() returned a non-nil error value,
	// i.e. the code ran fine but the action it performed failed.
	ToolError string
	Duration  time.Duration
}

// timeout bounds how long a single generated snippet is allowed to run.
// yaegi cannot forcibly preempt a runaway goroutine, so this is a soft
// deadline: the caller gets control back, but a truly stuck goroutine (e.g.
// `for {}`) leaks until the process exits. Good enough for a PoC; a
// production system would run untrusted code in a subprocess instead.
const timeout = 10 * time.Second

// Executor evaluates generated code against the fixed tools API.
type Executor struct {
	exports interp.Exports
}

// New builds an Executor. It must be called after tools.SetWorkDir so the
// tools package is sandboxed to the intended directory before any code runs.
func New() *Executor {
	return &Executor{exports: toolExports()}
}

// Run compiles and executes source, which must be a self-contained Go file
// declaring `package main`, importing "tools" as needed, and defining
//
//	func Run() (string, error)
//
// as its entry point. Run() is invoked automatically after the source is
// loaded, and its return values are captured in Result.
func (e *Executor) Run(ctx context.Context, source string) (Result, error) {
	start := time.Now()

	var stdout bytes.Buffer
	i := interp.New(interp.Options{Stdout: &stdout, Stderr: &stdout})
	if err := i.Use(allowedStdlib()); err != nil {
		return Result{}, fmt.Errorf("internal: loading stdlib symbols: %w", err)
	}
	if err := i.Use(e.exports); err != nil {
		return Result{}, fmt.Errorf("internal: loading tools symbols: %w", err)
	}

	type evalOutcome struct {
		res reflect.Value
		err reflect.Value
	}
	done := make(chan evalOutcome, 1)
	fail := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fail <- fmt.Errorf("panic during execution: %v", r)
			}
		}()
		if _, err := i.Eval(source); err != nil {
			fail <- fmt.Errorf("compile error: %w", err)
			return
		}
		// Run() may return (string, error); yaegi's Eval only surfaces the
		// first return value of a call expression, so we bind both to
		// package-level vars first and read them back individually.
		if _, err := i.Eval("__codeactRes, __codeactErr := Run()"); err != nil {
			fail <- fmt.Errorf("Run() invocation error: %w", err)
			return
		}
		resVal, err := i.Eval("__codeactRes")
		if err != nil {
			fail <- fmt.Errorf("internal: reading result: %w", err)
			return
		}
		errVal, err := i.Eval("__codeactErr")
		if err != nil {
			fail <- fmt.Errorf("internal: reading error result: %w", err)
			return
		}
		done <- evalOutcome{res: resVal, err: errVal}
	}()

	select {
	case <-ctx.Done():
		return Result{}, fmt.Errorf("execution cancelled: %w", ctx.Err())
	case <-time.After(timeout):
		return Result{}, fmt.Errorf("execution timed out after %s", timeout)
	case err := <-fail:
		return Result{Stdout: stdout.String(), Duration: time.Since(start)}, err
	case outcome := <-done:
		result := Result{
			Stdout:      stdout.String(),
			ReturnValue: fmt.Sprintf("%v", outcome.res.Interface()),
			Duration:    time.Since(start),
		}
		if outcome.err.IsValid() && !outcome.err.IsNil() {
			if goErr, ok := outcome.err.Interface().(error); ok {
				result.ToolError = goErr.Error()
			}
		}
		return result, nil
	}
}

// allowedStdlibPackages is the set of standard library packages exposed to
// generated code. All filesystem, process, and network access must go
// through the tools package instead, so packages like os, io, os/exec,
// net, and syscall are deliberately excluded: without them, generated code
// has no way to escape the tools.SetWorkDir sandbox even if it imports the
// package directly instead of calling tools.* functions.
var allowedStdlibPackages = map[string]bool{
	"fmt/fmt":            true,
	"strings/strings":    true,
	"strconv/strconv":    true,
	"errors/errors":      true,
	"sort/sort":          true,
	"time/time":          true,
	"regexp/regexp":      true,
	"bytes/bytes":        true,
	"unicode/unicode":    true,
	"unicode/utf8/utf8":  true,
	"math/math":          true,
	"encoding/json/json": true,
}

// allowedStdlib returns a filtered copy of yaegi's full stdlib symbol table
// containing only allowedStdlibPackages.
func allowedStdlib() interp.Exports {
	filtered := make(interp.Exports, len(allowedStdlibPackages))
	for pkg, symbols := range stdlib.Symbols {
		if allowedStdlibPackages[pkg] {
			filtered[pkg] = symbols
		}
	}
	return filtered
}

// toolExports builds the yaegi symbol table for the "tools" package by
// reflecting over the exported functions in internal/tools.
func toolExports() interp.Exports {
	return interp.Exports{
		"tools/tools": {
			"ReadFile":        reflect.ValueOf(tools.ReadFile),
			"WriteFile":       reflect.ValueOf(tools.WriteFile),
			"ListDir":         reflect.ValueOf(tools.ListDir),
			"FileTree":        reflect.ValueOf(tools.FileTree),
			"Grep":            reflect.ValueOf(tools.Grep),
			"CountLinesByExt": reflect.ValueOf(tools.CountLinesByExt),
			"RunCommand":      reflect.ValueOf(tools.RunCommand),
		},
	}
}
