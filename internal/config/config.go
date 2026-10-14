package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultCacheTTL = 5 * time.Minute

type Config struct {
	Roots       []string `yaml:"roots"`
	Extras      []string `yaml:"extras"`
	CacheTTLStr string   `yaml:"cache_ttl"`
}

func (c *Config) CacheTTL() time.Duration {
	if c.CacheTTLStr == "" {
		return DefaultCacheTTL
	}
	d, err := time.ParseDuration(c.CacheTTLStr)
	if err != nil {
		return DefaultCacheTTL
	}
	return d
}

func DefaultPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "lsproj", "config.yaml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "lsproj", "config.yaml")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) ExpandRoots() []string {
	return expandPaths(c.Roots)
}

func (c *Config) ExpandExtras() []string {
	return expandPaths(c.Extras)
}

func expandPaths(paths []string) []string {
	var expanded []string
	for _, p := range paths {
		if len(p) > 0 && p[0] == '~' {
			home, _ := os.UserHomeDir()
			p = filepath.Join(home, p[1:])
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		expanded = append(expanded, abs)
	}
	return expanded
}
