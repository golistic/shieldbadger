// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type ShieldBadger struct {
	config           *Config
	FeedbackCallback func(format string, a ...any)
}

func NewShieldBadger(configFiles ...string) (*ShieldBadger, error) {
	if len(configFiles) > 1 {
		panic("only one configuration file is supported")
	}
	sb := &ShieldBadger{}

	if configFiles[0] != "" {
		var err error
		sb.config, err = newConfigFromFile(configFiles[0])
		if err != nil {
			return nil, err
		}
	}

	return sb, nil
}

func (sb *ShieldBadger) Fetch() error {
	if err := os.Mkdir(sb.config.DestFolder, 0770); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("badges: making folder %s (%w)", sb.config.DestFolder, err)
		}
	}

	for _, badge := range sb.config.Badges {
		if err := sb.fetchBadge(badge); err != nil {
			return err
		}

		sb.FeedbackCallback("stored badge %s\n", badge.Name)
	}

	return nil
}

func (sb *ShieldBadger) fetchBadge(badge *Badge) error {
	baseURL, err := url.Parse(sb.config.URLShieldsIO)
	if err != nil {
		return fmt.Errorf("%s: parsing URL %s (%w)", badge.Name, sb.config.URLShieldsIO, err)
	}

	u := *baseURL

	if values, err := badge.urlQueryParameters(sb.config); err != nil {
		return fmt.Errorf("%s: preparing URL query parameters %w", badge.Name, err)
	} else {
		u.RawQuery = values.Encode()
	}

	client := http.DefaultClient
	resp, err := client.Get(u.String())
	if err != nil {
		return fmt.Errorf("badges: %s: getting URL %s (%w)", badge.Name, u.String(), err)
	}
	defer func() { _ = resp.Body.Close() }()

	f := filepath.Join(sb.config.DestFolder, badge.Name+".svg")
	fp, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return fmt.Errorf("badges: %s: opening file %s (%w)", badge.Name, f, err)
	}

	if _, err := io.Copy(fp, resp.Body); err != nil {
		return fmt.Errorf("badges: %s: writing to file %s (%w)", badge.Name, f, err)
	}

	if err := fp.Close(); err != nil {
		return fmt.Errorf("badges: %s: closing file %s (%w)", badge.Name, f, err)
	}

	return nil
}

func (sb *ShieldBadger) feedback(format string, a ...any) {
	if sb.FeedbackCallback != nil {
		sb.FeedbackCallback(format, a...)
	}
}
