package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	envDatabase   = "BOOKS_DB"
	configDirName = "books"
	configFile    = "config.toml"
)

type Source string

const (
	SourceEnv        Source = "env"
	SourceConfigFile Source = "config_file"
	SourceDefault    Source = "default"
)

type Config struct {
	DatabasePath string
	ConfigPath   string
	ConfigExists bool
	Source       Source
}

type fileConfig struct {
	Database string `toml:"database"`
}

func Resolve() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home directory: %w", err)
	}

	cfgPath := filepath.Join(home, ".config", configDirName, configFile)

	if v := os.Getenv(envDatabase); v != "" {
		return Config{
			DatabasePath: v,
			ConfigPath:   cfgPath,
			ConfigExists: fileExists(cfgPath),
			Source:       SourceEnv,
		}, nil
	}

	if fileExists(cfgPath) {
		var fc fileConfig
		if _, err := toml.DecodeFile(cfgPath, &fc); err != nil {
			return Config{}, fmt.Errorf("read config file %s: %w", cfgPath, err)
		}
		if fc.Database != "" {
			return Config{
				DatabasePath: fc.Database,
				ConfigPath:   cfgPath,
				ConfigExists: true,
				Source:       SourceConfigFile,
			}, nil
		}
	}

	defaultPath := filepath.Join(home, ".local", "share", configDirName, "books.db")
	return Config{
		DatabasePath: defaultPath,
		ConfigPath:   cfgPath,
		ConfigExists: fileExists(cfgPath),
		Source:       SourceDefault,
	}, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
