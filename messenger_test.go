// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"testing"

	"github.com/golistic/xt"
)

func TestMessenger(t *testing.T) {
	t.Run("not registered messenger", func(t *testing.T) {
		_, err := NewMessenger("something.not.available", nil)
		xt.KO(t, err)
		xt.Eq(t, ErrMessengerNotFound, err)
	})

	t.Run("all registered messengers return no error", func(t *testing.T) {
		for msger := range registry {
			_, err := NewMessenger(msger, nil)
			xt.OK(t, err)
		}
	})
}
