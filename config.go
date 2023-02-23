// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Style        string                     `json:"style"`
	URLShieldsIO string                     `json:"urlShieldsIO"`
	DestFolder   string                     `json:"destFolder"`
	Badges       []*Badge                   `json:"badges"`
	Messengers   map[string]json.RawMessage `json:"messengers"`
}

func DefaultConfig() *Config {
	return &Config{
		Style:        defaultStyle,
		URLShieldsIO: "https://img.shields.io/static/v1",
		DestFolder:   "_badges",
	}
}

func newConfigFromFile(name string) (*Config, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening configuration file %s (%w)", name, err)
	}

	c, err := newConfigFromReader(fp)
	if err != nil {
		return nil, fmt.Errorf("reading configuration file %s (%w)", name, err)
	}

	return c, nil
}

func newConfigFromReader(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
