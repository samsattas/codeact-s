// Command web is a browser front-end for the same CodeAct agent core used
// by cmd/agent: it accepts a task over HTTP and streams each
// generate/execute/fix attempt back to the page as it happens, so you can
// watch the model write Go code, watch it run, and watch it self-correct.
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
// the browser for a run: a "start" (the resolved sandbox directory for this
// run), a "step" (one generate/execute attempt), or a terminal
// "done"/"fatal" message.
type streamMessage struct {
	Type string `json:"type"`

	// start fields
	Dir string `json:"dir,omitempty"`

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

	// Fail fast if the default workdir is bad, but the resulting Sandbox
	// isn't kept: every request builds its own (defaulting to this
	// directory, or overriding it with its own "dir"), so a bad directory
	// picked from the UI later can't wedge the whole server.
	defaultSandbox, err := tools.NewSandbox(*workDir)
	if err != nil {
		log.Fatalf("invalid -workdir: %v", err)
	}
	provider, err := llm.FromEnv()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("using model provider: %s", provider.Name())
	exec := executor.New()

	srv := &server{provider: provider, exec: exec, maxAttempts: *maxAttempts, defaultDir: defaultSandbox.Root()}

	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("embedding static assets: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticContent)))
	mux.HandleFunc("/api/info", srv.handleInfo)
	mux.HandleFunc("/api/run", srv.handleRun)
	mux.HandleFunc("/api/browse", srv.handleBrowse)

	log.Printf("listening on http://localhost%s (default sandbox: %s)", *addr, defaultSandbox.Root())
	log.Fatal(http.ListenAndServe(*addr, mux))
}

type server struct {
	provider    llm.Provider
	exec        *executor.Executor
	maxAttempts int
	// defaultDir is used for a run when the request doesn't specify its own
	// "dir" — it is not a fixed sandbox for the process, just a fallback.
	defaultDir string
}

func (s *server) handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"provider":    s.provider.Name(),
		"workdir":     s.defaultDir,
		"maxAttempts": s.maxAttempts,
	})
}

// browseEntry is one subdirectory returned by handleBrowse.
type browseEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// handleBrowse lists the subdirectories of ?path= (defaulting to the
// server's default dir) so the web UI's folder picker can navigate the
// server's filesystem without the browser ever seeing (or needing) a real
// absolute path from a native file input, which browsers deliberately
// don't expose for security reasons.
func (s *server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("path")
	if dir == "" {
		dir = s.defaultDir
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !info.IsDir() {
		http.Error(w, fmt.Sprintf("%q is not a directory", abs), http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dirs := make([]browseEntry, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || strings.HasPrefix(name, ".") {
			continue
		}
		dirs = append(dirs, browseEntry{Name: name, Path: filepath.Join(abs, name)})
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name < dirs[j].Name })

	parent := filepath.Dir(abs)
	if parent == abs {
		parent = "" // already at the filesystem root
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"path":    abs,
		"parent":  parent,
		"entries": dirs,
	})
}

func (s *server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Task        string `json:"task"`
		Dir         string `json:"dir"`
		MaxAttempts int    `json:"maxAttempts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Task == "" {
		http.Error(w, "expected JSON body with a non-empty \"task\"", http.StatusBadRequest)
		return
	}

	dir := req.Dir
	if dir == "" {
		dir = s.defaultDir
	}
	sandbox, err := tools.NewSandbox(dir)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid \"dir\": %v", err), http.StatusBadRequest)
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

	send(streamMessage{Type: "start", Dir: sandbox.Root()})

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

	ag := agent.New(s.provider, s.exec, sandbox, maxAttempts, onStep)
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
