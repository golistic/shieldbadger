// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import "net/url"

const defaultStyle = "flat"

var supportedStyles = map[string]struct{}{
	"flat":          {},
	"plastic":       {},
	"flat-square":   {},
	"for-the-badge": {},
	"social":        {},
}

// Badge contains information on the content and the style of
// a badge produced by Shields.io.
type Badge struct {
	Label       string `json:"label"`
	Message     string `json:"message"`
	MessageFunc string `json:"messageFunc"`
	Color       string `json:"color"`
	Style       string `json:"style"`
	Name        string `json:"name"`

	config *Config
}

// urlQueryParameters produces the query parameters of the URL used to communicated
// with Shields.io. The `url.Values` must be encoded before used.
func (b Badge) urlQueryParameters(cfg *Config) (*url.Values, error) {
	v := &url.Values{}
	v.Set("label", b.Label)
	v.Set("color", b.Color)

	if _, ok := supportedStyles[b.Style]; ok {
		v.Set("style", b.Style)
	} else {
		v.Set("style", cfg.Style)
	}

	if b.MessageFunc != "" {
		mr, err := NewMessenger(b.MessageFunc, b.config)
		if err != nil {
			return nil, err
		}

		m, err := mr.Message()
		if err != nil {
			return nil, err
		}

		v.Set("message", m)
	} else {
		v.Set("message", b.Message)
	}

	return v, nil
}
