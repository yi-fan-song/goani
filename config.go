/**
 * Copyright (C) 2021 Yi Fan Song <yfsong00@gmail.com>
 *
 * This file is part of Goani.
 *
 * Goani is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Goani is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Goani.  If not, see <https://www.gnu.org/licenses/>.
 **/

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
