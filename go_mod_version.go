// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/mod/modfile"
)

var cfgKeyGoModVersion = "go.mod.version"

func readGoMod(path string) (*modfile.File, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening go.mod (%w)", err)
	}
	data, err := io.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("reading go.mod (%w)", err)
	}
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing go.mod (%w)", err)
	}

	return modFile, err
}

type GoModVersion struct {
	Path string `json:"path"`
}

var _ Messenger = (*GoModVersion)(nil)

func NewGoModVersion(cfg *Config) (Messenger, error) {
	m := &GoModVersion{
		Path: "go.mod",
	}

	if cfg != nil {
		if r, ok := cfg.Messengers[cfgKeyGoModVersion]; ok {
			if err := json.Unmarshal(r, m); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m GoModVersion) Message() (string, error) {
	if m.Path == "" {
		m.Path = "go.mod"
	}

	modFile, err := readGoMod(m.Path)
	if err != nil {
		return "", err
	}

	return modFile.Go.Version, nil
}
