package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.IsRuleEnabled("no-unknown-diagram-type") {
		t.Error("expected no-unknown-diagram-type to be enabled by default")
	}
	if !cfg.IsRuleEnabled("valid-direction") {
		t.Error("expected valid-direction to be enabled by default")
	}
	if cfg.IsRuleEnabled("nonexistent-rule") {
		t.Error("expected nonexistent rule to be disabled")
	}
}

func TestRuleSeverity(t *testing.T) {
	cfg := DefaultConfig()
	if s := cfg.RuleSeverity("no-unknown-diagram-type"); s != SeverityError {
		t.Errorf("expected error severity, got %q", s)
	}
	if s := cfg.RuleSeverity("node-has-label"); s != SeverityInfo {
		t.Errorf("expected info severity, got %q", s)
	}
	if s := cfg.RuleSeverity("nonexistent"); s != SeverityWarning {
		t.Errorf("expected default warning severity, got %q", s)
	}
}

func TestLoadConfig_NonexistentFile(t *testing.T) {
	cfg, err := LoadConfig("/tmp/nonexistent-mermaid-lint-config.json")
	if err != nil {
		t.Fatal(err)
	}
	// Should return default config
	if !cfg.IsRuleEnabled("no-unknown-diagram-type") {
		t.Error("expected default config when file doesn't exist")
	}
}

func TestLoadConfig_CustomFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	data := []byte(`{"rules": {"no-unknown-diagram-type": {"enabled": false, "severity": "info"}}}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.IsRuleEnabled("no-unknown-diagram-type") {
		t.Error("expected no-unknown-diagram-type to be disabled via config")
	}
}
