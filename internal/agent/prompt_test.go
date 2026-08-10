package agent

import (
	"strings"
	"testing"
)

func TestDiagnosisHintUndefined(t *testing.T) {
	hint := diagnosisHint("16:9: undefined: fmt")
	if hint == "" {
		t.Fatalf("expected a hint for an undefined-symbol error")
	}
	if !strings.Contains(hint, `"fmt"`) {
		t.Fatalf("expected hint to name the undefined symbol, got: %q", hint)
	}
}

func TestDiagnosisHintUnusedImport(t *testing.T) {
	hint := diagnosisHint(`5:2: "strings" imported and not used`)
	if hint == "" {
		t.Fatalf("expected a hint for an unused-import error")
	}
}

func TestDiagnosisHintUnusedVar(t *testing.T) {
	hint := diagnosisHint("7:2: x declared and not used")
	if hint == "" {
		t.Fatalf("expected a hint for an unused-variable error")
	}
}

func TestDiagnosisHintExitStatus(t *testing.T) {
	hint := diagnosisHint("Run() returned an error: go vet failed: exit status 1")
	if hint == "" {
		t.Fatalf("expected a hint for a non-zero exit status error")
	}
	if !strings.Contains(hint, "RunCommand") {
		t.Fatalf("expected hint to reference RunCommand, got: %q", hint)
	}
}

func TestDiagnosisHintUnknownErrorIsEmpty(t *testing.T) {
	if hint := diagnosisHint("something completely different"); hint != "" {
		t.Fatalf("expected no hint for an unrecognized failure, got: %q", hint)
	}
}

func TestRetryPromptIncludesHintAndFailure(t *testing.T) {
	prompt := retryPrompt("do something", "package main\n", "16:9: undefined: fmt")
	if !strings.Contains(prompt, "16:9: undefined: fmt") {
		t.Fatalf("expected retry prompt to include the raw failure message")
	}
	if !strings.Contains(prompt, `"fmt"`) {
		t.Fatalf("expected retry prompt to include the diagnosis hint, got: %q", prompt)
	}
}
