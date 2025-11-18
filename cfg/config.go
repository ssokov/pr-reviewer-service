package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"sslmode"`
}

type ServerConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	IsDevel bool   `toml:"is_devel"`
}

type APIKeysConfig struct {
	MaxBot string `toml:"max_bot"`
}

type Config struct {
	Database DBConfig      `toml:"database"`
	Server   ServerConfig  `toml:"server"`
	APIKeys  APIKeysConfig `toml:"api_keys"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("max_superuser://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}
