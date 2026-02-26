// Package config handles configuration for the mermaid-lint tool.
package config

import (
	"encoding/json"
	"os"
)

// Severity represents how severe a lint finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// RuleConfig holds configuration for a single lint rule.
type RuleConfig struct {
	Enabled  bool     `json:"enabled"`
	Severity Severity `json:"severity"`
}

// Config holds the complete linter configuration.
type Config struct {
	Rules map[string]RuleConfig `json:"rules"`
}

// DefaultConfig returns the default configuration with all rules enabled.
func DefaultConfig() *Config {
	return &Config{
		Rules: map[string]RuleConfig{
			"no-unknown-diagram-type": {Enabled: true, Severity: SeverityError},
			"no-empty-diagram":        {Enabled: true, Severity: SeverityWarning},
			"valid-direction":          {Enabled: true, Severity: SeverityError},
			"no-duplicate-node-ids":    {Enabled: true, Severity: SeverityWarning},
			"node-has-label":           {Enabled: true, Severity: SeverityInfo},
			"no-orphan-nodes":          {Enabled: true, Severity: SeverityWarning},
		},
	}
}

// LoadConfig loads configuration from a JSON file.
// If the file does not exist, it returns the default configuration.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsRuleEnabled checks whether a rule is enabled in the config.
func (c *Config) IsRuleEnabled(name string) bool {
	rule, ok := c.Rules[name]
	if !ok {
		return false
	}
	return rule.Enabled
}

// RuleSeverity returns the severity for a rule, defaulting to warning.
func (c *Config) RuleSeverity(name string) Severity {
	rule, ok := c.Rules[name]
	if !ok {
		return SeverityWarning
	}
	return rule.Severity
}
