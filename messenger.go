// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import "fmt"

var ErrMessengerNotFound = fmt.Errorf("messenger not available")

type NewMessengerFunc func(cfg *Config) (Messenger, error)

var registry = map[string]NewMessengerFunc{
	"go.mod.version":      NewGoModVersion,
	"go.reportcard.grade": NewGoReportCardGrade,
}

type Messenger interface {
	Message() (string, error)
}

func NewMessenger(name string, cfg *Config) (Messenger, error) {
	f, ok := registry[name]
	if !ok {
		return nil, ErrMessengerNotFound
	}
	return f(cfg)
}
