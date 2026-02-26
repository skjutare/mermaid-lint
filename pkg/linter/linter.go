// Package linter provides the core linting engine for Mermaid diagrams.
package linter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/skjutare/mermaid-lint/pkg/config"
	"github.com/skjutare/mermaid-lint/pkg/parser"
)

// Finding represents a single lint finding.
type Finding struct {
	Rule     string
	Severity config.Severity
	Message  string
	File     string
	Line     int
}

// String returns a human-readable representation of the finding.
func (f Finding) String() string {
	loc := f.File
	if f.Line > 0 {
		loc = fmt.Sprintf("%s:%d", f.File, f.Line)
	}
	return fmt.Sprintf("%s [%s] %s (%s)", loc, f.Severity, f.Message, f.Rule)
}

// Linter runs lint rules against Mermaid diagrams.
type Linter struct {
	Config *config.Config
	Rules  []Rule
}

// New creates a new Linter with the given configuration.
func New(cfg *config.Config) *Linter {
	l := &Linter{Config: cfg}
	l.Rules = AllRules()
	return l
}

// LintFile lints a single file (either .mmd or .md).
func (l *Linter) LintFile(path string) ([]Finding, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".mmd", ".mermaid":
		return l.lintMermaidFile(path)
	case ".md", ".markdown":
		return l.lintMarkdownFile(path)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func (l *Linter) lintMermaidFile(path string) ([]Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source := string(data)
	diagram := parser.Parse(source, 1)
	return l.lintDiagram(diagram, path), nil
}

func (l *Linter) lintMarkdownFile(path string) ([]Finding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return l.LintMarkdownReader(f, path)
}

// LintMarkdownReader lints mermaid blocks extracted from a markdown reader.
func (l *Linter) LintMarkdownReader(r io.Reader, filename string) ([]Finding, error) {
	blocks, err := parser.ExtractMermaidBlocks(r)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	for _, block := range blocks {
		// StartLine+1 because the mermaid source starts on the line after the fence
		diagram := parser.Parse(block.Source, block.StartLine+1)
		findings = append(findings, l.lintDiagram(diagram, filename)...)
	}
	return findings, nil
}

// LintSource lints raw mermaid source code.
func (l *Linter) LintSource(source string, filename string) []Finding {
	diagram := parser.Parse(source, 1)
	return l.lintDiagram(diagram, filename)
}

func (l *Linter) lintDiagram(d *parser.Diagram, filename string) []Finding {
	var findings []Finding
	for _, rule := range l.Rules {
		if !l.Config.IsRuleEnabled(rule.Name()) {
			continue
		}
		severity := l.Config.RuleSeverity(rule.Name())
		for _, f := range rule.Check(d) {
			f.File = filename
			f.Severity = severity
			findings = append(findings, f)
		}
	}
	return findings
}
