package linter

import (
	"strings"
	"testing"

	"github.com/skjutare/mermaid-lint/pkg/config"
)

func TestLintSource_UnknownDiagramType(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	findings := l.LintSource("badtype\n  stuff", "test.mmd")
	found := findByRule(findings, "no-unknown-diagram-type")
	if len(found) == 0 {
		t.Error("expected finding for unknown diagram type")
	}
}

func TestLintSource_ValidFlowchart(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	findings := l.LintSource("flowchart LR\n  A[Start] --> B[End]", "test.mmd")
	errors := findBySeverity(findings, config.SeverityError)
	if len(errors) > 0 {
		t.Errorf("expected no errors for valid flowchart, got: %v", errors)
	}
}

func TestLintSource_InvalidDirection(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	findings := l.LintSource("flowchart XX\n  A --> B", "test.mmd")
	found := findByRule(findings, "valid-direction")
	if len(found) == 0 {
		t.Error("expected finding for invalid direction")
	}
}

func TestLintSource_EmptyDiagram(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	findings := l.LintSource("flowchart LR\n", "test.mmd")
	found := findByRule(findings, "no-empty-diagram")
	if len(found) == 0 {
		t.Error("expected finding for empty diagram")
	}
}

func TestLintSource_NodeHasLabel(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	findings := l.LintSource("flowchart LR\n  A --> B", "test.mmd")
	found := findByRule(findings, "node-has-label")
	if len(found) != 2 {
		t.Errorf("expected 2 node-has-label findings, got %d", len(found))
	}
}

func TestLintSource_OrphanNodes(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	// C is orphaned - not connected to any edge
	findings := l.LintSource("flowchart LR\n  A --> B\n  C[Orphan]", "test.mmd")
	found := findByRule(findings, "no-orphan-nodes")
	if len(found) == 0 {
		t.Error("expected finding for orphan node")
	}
}

func TestLintSource_DisabledRule(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Rules["node-has-label"] = config.RuleConfig{Enabled: false, Severity: config.SeverityInfo}
	l := New(cfg)

	findings := l.LintSource("flowchart LR\n  A --> B", "test.mmd")
	found := findByRule(findings, "node-has-label")
	if len(found) != 0 {
		t.Error("expected no findings for disabled rule")
	}
}

func TestLintMarkdownReader(t *testing.T) {
	cfg := config.DefaultConfig()
	l := New(cfg)

	md := "# Doc\n\n```mermaid\nflowchart XX\n  A --> B\n```\n"
	findings, err := l.LintMarkdownReader(strings.NewReader(md), "test.md")
	if err != nil {
		t.Fatal(err)
	}
	found := findByRule(findings, "valid-direction")
	if len(found) == 0 {
		t.Error("expected finding for invalid direction in markdown")
	}
	if found[0].File != "test.md" {
		t.Errorf("file = %q, want %q", found[0].File, "test.md")
	}
}

func TestFindingString(t *testing.T) {
	f := Finding{
		Rule:     "test-rule",
		Severity: config.SeverityError,
		Message:  "test message",
		File:     "test.mmd",
		Line:     5,
	}
	s := f.String()
	if !strings.Contains(s, "test.mmd:5") {
		t.Errorf("expected file:line in output, got %q", s)
	}
	if !strings.Contains(s, "error") {
		t.Errorf("expected severity in output, got %q", s)
	}
}

func findByRule(findings []Finding, rule string) []Finding {
	var result []Finding
	for _, f := range findings {
		if f.Rule == rule {
			result = append(result, f)
		}
	}
	return result
}

func findBySeverity(findings []Finding, sev config.Severity) []Finding {
	var result []Finding
	for _, f := range findings {
		if f.Severity == sev {
			result = append(result, f)
		}
	}
	return result
}
