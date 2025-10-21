package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration file structure
type Config struct {
	Defaults ScanDefaults           `yaml:"defaults"`
	Profiles map[string]ScanProfile `yaml:"profiles"`
}

// ScanDefaults holds default scan settings
type ScanDefaults struct {
	Path   string `yaml:"path"`
	Output string `yaml:"output"`
	NoWarn bool   `yaml:"no_warn"`
	Filter string `yaml:"filter"`
	Key    string `yaml:"key"`
	Value  string `yaml:"value"`
}

// ScanProfile represents a named configuration profile
type ScanProfile struct {
	Path        string `yaml:"path"`
	Output      string `yaml:"output"`
	NoWarn      bool   `yaml:"no_warn"`
	Filter      string `yaml:"filter"`
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

// LoadConfig loads the config file from ~/.konfetti.yaml
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".konfetti.yaml")

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Defaults: ScanDefaults{Output: "text"},
			Profiles: make(map[string]ScanProfile),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set default output if not specified
	if cfg.Defaults.Output == "" {
		cfg.Defaults.Output = "text"
	}

	return &cfg, nil
}

// GetSampleConfig returns a sample configuration file content
func GetSampleConfig() string {
	return `# Konfetti Configuration File
# Save this as ~/.konfetti.yaml

# Default settings applied to all scans (can be overridden by CLI flags)
defaults:
  output: text        # Default output format: text, json, table
  no_warn: false      # Suppress warning messages
  # path: /etc        # Default scan path (omit to use current directory)
  # filter: ""        # Default filename filter
  # key: ""           # Default key filter
  # value: ""         # Default value filter

# Named profiles for common scanning scenarios
profiles:
  # Example: Scan for debug settings across all configs
  debug:
    description: "Find all debug-related settings"
    key: debug
    output: table
    no_warn: true

  # Example: Find production configs
  prod:
    description: "Find production environment settings"
    value: prod
    output: json
    no_warn: true

  # Example: Scan system-wide configs
  system:
    description: "Scan system configuration directories"
    path: /etc
    output: text
    no_warn: false

  # Example: Docker configs
  docker:
    description: "Find Docker-related configurations"
    filter: docker
    output: table

  # Example: Security audit - find sensitive keys
  security:
    description: "Find potentially sensitive configuration keys"
    key: password|secret|key|token|credential
    output: table
    no_warn: true

# Usage:
#   konfetti scan                    # Uses defaults
#   konfetti scan --profile debug    # Uses debug profile
#   konfetti scan --profile prod --output table  # Profile + CLI override
`
}
