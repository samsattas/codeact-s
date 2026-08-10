// Command web is a browser front-end for the same CodeAct agent core used
// by cmd/agent: it accepts a task over HTTP and streams each
// generate/execute/fix attempt back to the page as it happens, so you can
// watch the model write Go code, watch it run, and watch it self-correct.
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"codeact-agent/internal/agent"
	"codeact-agent/internal/executor"
	"codeact-agent/internal/llm"
	"codeact-agent/internal/tools"
)

//go:embed static
var staticFS embed.FS

// maxAttemptsCap bounds the per-request maxAttempts override from the web
// UI, so a client can't force an unbounded number of model calls.
const maxAttemptsCap = 10

// streamMessage is one line of the newline-delimited JSON stream sent to
// the browser for a run: either a "step" (one generate/execute attempt) or
// a terminal "done"/"fatal" message.
type streamMessage struct {
	Type string `json:"type"`

	// step fields
	Attempt   int    `json:"attempt,omitempty"`
	Code      string `json:"code,omitempty"`
	Stdout    string `json:"stdout,omitempty"`
	ToolError string `json:"toolError,omitempty"`
	ExecError string `json:"execError,omitempty"`

	// done/fatal fields
	Success     bool   `json:"success,omitempty"`
	FinalAnswer string `json:"finalAnswer,omitempty"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
}

func main() {
	var (
		addr        = flag.String("addr", ":8080", "address to listen on")
		workDir     = flag.String("workdir", ".", "directory the agent is allowed to operate on")
		maxAttempts = flag.Int("max-attempts", 3, "max generate/execute/fix attempts per task")
	)
	flag.Parse()

	if err := tools.SetWorkDir(*workDir); err != nil {
		log.Fatalf("invalid -workdir: %v", err)
	}
	provider, err := llm.FromEnv()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("using model provider: %s", provider.Name())
	exec := executor.New()

	srv := &server{provider: provider, exec: exec, maxAttempts: *maxAttempts, workDir: tools.WorkDir()}

	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("embedding static assets: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticContent)))
	mux.HandleFunc("/api/info", srv.handleInfo)
	mux.HandleFunc("/api/run", srv.handleRun)

	log.Printf("listening on http://localhost%s (sandbox: %s)", *addr, tools.WorkDir())
	log.Fatal(http.ListenAndServe(*addr, mux))
}

type server struct {
	provider    llm.Provider
	exec        *executor.Executor
	maxAttempts int
	workDir     string
}

func (s *server) handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"provider":    s.provider.Name(),
		"workdir":     s.workDir,
		"maxAttempts": s.maxAttempts,
	})
}

func (s *server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Task        string `json:"task"`
		MaxAttempts int    `json:"maxAttempts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Task == "" {
		http.Error(w, "expected JSON body with a non-empty \"task\"", http.StatusBadRequest)
		return
	}

	maxAttempts := s.maxAttempts
	if req.MaxAttempts > 0 {
		maxAttempts = req.MaxAttempts
		if maxAttempts > maxAttemptsCap {
			maxAttempts = maxAttemptsCap
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	send := func(m streamMessage) {
		_ = enc.Encode(m)
		flusher.Flush()
	}

	onStep := func(step agent.Step) {
		msg := streamMessage{
			Type:    "step",
			Attempt: step.Attempt,
			Code:    step.Code,
			Stdout:  step.Output.Stdout,
		}
		if step.ExecErr != nil {
			msg.ExecError = step.ExecErr.Error()
		}
		if step.Output.ToolError != "" {
			msg.ToolError = step.Output.ToolError
		}
		send(msg)
	}

	ag := agent.New(s.provider, s.exec, maxAttempts, onStep)
	outcome, err := ag.Do(r.Context(), req.Task)
	if err != nil {
		send(streamMessage{Type: "fatal", Message: err.Error()})
		return
	}

	done := streamMessage{Type: "done", Success: outcome.Success, FinalAnswer: outcome.FinalAnswer}
	if !outcome.Success && len(outcome.Steps) > 0 {
		last := outcome.Steps[len(outcome.Steps)-1]
		if last.ExecErr != nil {
			done.Error = last.ExecErr.Error()
		} else {
			done.Error = last.Output.ToolError
		}
	}
	send(done)
}
