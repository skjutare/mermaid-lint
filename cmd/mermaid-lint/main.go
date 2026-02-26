// Command mermaid-lint is a linter for Mermaid diagram files (.mmd)
// and Mermaid code blocks embedded in Markdown files (.md).
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skjutare/mermaid-lint/pkg/config"
	"github.com/skjutare/mermaid-lint/pkg/linter"
)

var version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	configPath := flag.String("config", ".mermaid-lint.json", "path to configuration file")
	listRules := flag.Bool("list-rules", false, "list all available lint rules")
	showVersion := flag.Bool("version", false, "show version")
	severityFilter := flag.String("severity", "", "only show findings at this severity or above (info, warning, error)")
	outputFormat := flag.String("format", "text", "output format: text or json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mermaid-lint [flags] <files or directories...>\n\n")
		fmt.Fprintf(os.Stderr, "A linter for Mermaid diagram files (.mmd) and Mermaid blocks in Markdown (.md).\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  mermaid-lint diagram.mmd\n")
		fmt.Fprintf(os.Stderr, "  mermaid-lint docs/\n")
		fmt.Fprintf(os.Stderr, "  mermaid-lint --severity warning *.md\n")
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("mermaid-lint %s\n", version)
		return 0
	}

	if *listRules {
		printRules()
		return 0
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: no files or directories specified")
		flag.Usage()
		return 1
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		return 1
	}

	l := linter.New(cfg)

	files, err := collectFiles(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no supported files found")
		return 1
	}

	var allFindings []linter.Finding
	for _, file := range files {
		findings, err := l.LintFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error linting %s: %v\n", file, err)
			continue
		}
		allFindings = append(allFindings, findings...)
	}

	// Apply severity filter
	if *severityFilter != "" {
		allFindings = filterBySeverity(allFindings, config.Severity(*severityFilter))
	}

	printFindings(allFindings, *outputFormat)

	// Count errors and warnings for exit code
	errorCount := 0
	for _, f := range allFindings {
		if f.Severity == config.SeverityError {
			errorCount++
		}
	}

	if errorCount > 0 {
		return 1
	}
	return 0
}

func printRules() {
	rules := linter.AllRules()
	cfg := config.DefaultConfig()
	fmt.Println("Available rules:")
	fmt.Println()
	for _, rule := range rules {
		status := "disabled"
		if cfg.IsRuleEnabled(rule.Name()) {
			status = string(cfg.RuleSeverity(rule.Name()))
		}
		fmt.Printf("  %-30s %s [%s]\n", rule.Name(), rule.Description(), status)
	}
}

func printFindings(findings []linter.Finding, format string) {
	if format == "json" {
		printFindingsJSON(findings)
		return
	}
	for _, f := range findings {
		fmt.Println(f.String())
	}
	if len(findings) > 0 {
		fmt.Printf("\n%d finding(s)\n", len(findings))
	}
}

func printFindingsJSON(findings []linter.Finding) {
	fmt.Println("[")
	for i, f := range findings {
		comma := ","
		if i == len(findings)-1 {
			comma = ""
		}
		fmt.Printf("  {\"rule\": %q, \"severity\": %q, \"message\": %q, \"file\": %q, \"line\": %d}%s\n",
			f.Rule, f.Severity, f.Message, f.File, f.Line, comma)
	}
	fmt.Println("]")
}

func filterBySeverity(findings []linter.Finding, minSeverity config.Severity) []linter.Finding {
	order := map[config.Severity]int{
		config.SeverityInfo:    0,
		config.SeverityWarning: 1,
		config.SeverityError:   2,
	}
	minOrder, ok := order[minSeverity]
	if !ok {
		return findings
	}
	var filtered []linter.Finding
	for _, f := range findings {
		if order[f.Severity] >= minOrder {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func collectFiles(args []string) ([]string, error) {
	var files []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("cannot access %s: %w", arg, err)
		}
		if info.IsDir() {
			dirFiles, err := walkDir(arg)
			if err != nil {
				return nil, err
			}
			files = append(files, dirFiles...)
		} else {
			if isSupportedFile(arg) {
				files = append(files, arg)
			}
		}
	}
	return files, nil
}

func walkDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := filepath.Base(path)
			// Skip hidden directories and common non-content directories
			if strings.HasPrefix(base, ".") || base == "node_modules" || base == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if isSupportedFile(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func isSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mmd" || ext == ".mermaid" || ext == ".md" || ext == ".markdown"
}
