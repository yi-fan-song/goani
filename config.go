package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Port    int
	Folders []string
}

const ConfigFile = "config.toml"

func LoadConfig() (*Config, error) {
	f, err := os.Open(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	var cfg Config
	err = toml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	return &cfg, nil
}

func CreateConfigIfNotExist() error {
	_, err := os.Stat(ConfigFile)
	if err == nil {
		// file exists
		return nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("unknown error while creating config: %w", err)
	}

	f, err := os.Create(ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	cfg := Config{
		Port:    8080,
		Folders: []string{"."},
	}

	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	f.Write(b)

	return nil
}

func LoadOrCreateConfig() (*Config, error) {
	err := CreateConfigIfNotExist()
	if err != nil {
		return nil, err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
