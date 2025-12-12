package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
)

// Config mirrors the legacy db.yaml structure to ease migration.
type Config struct {
	DBs map[string]*DBConfig `yaml:"dbs"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// LoadConfig loads database config from .db.yaml or db.yaml (compatible with legacy layout).
func LoadConfig() (*Config, error) {
	paths := []string{
		Tea.ConfigFile(".db.yaml"),
		Tea.ConfigFile("db.yaml"),
		filepath.Join("configs", ".db.yaml"),
		filepath.Join("configs", "db.yaml"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal %s: %w", p, err)
		}
		if len(cfg.DBs) == 0 {
			continue
		}
		return &cfg, nil
	}
	return nil, fmt.Errorf("db config not found (tried .db.yaml/db.yaml)")
}
