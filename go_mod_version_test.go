// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golistic/xt"
)

func TestNewGoModVersion(t *testing.T) {
	t.Run("this package its version in go mod", func(t *testing.T) {
		// whenever we change in go.mod, we must change this test; this is OK
		exp := "1.20"
		m, err := NewGoModVersion(nil)
		xt.OK(t, err)
		have, err := m.Message()
		xt.OK(t, err)
		xt.Eq(t, exp, have)
	})

	t.Run("go.mod is not available", func(t *testing.T) {
		m, err := NewGoModVersion(nil)
		xt.OK(t, err)
		gv := m.(*GoModVersion)
		gv.Path = "something/not/go.mod"
		_, err = m.Message()
		xt.KO(t, err)
		xt.Eq(t, "opening go.mod (open something/not/go.mod: no such file or directory)", err.Error())
	})

	t.Run("configured", func(t *testing.T) {
		cfg, err := newConfigFromReader(bytes.NewBufferString(fmt.Sprintf(`
{
	"style": "plastic",
	"messengers": {
		"%s": {
			"path": "my/go.mod"
		}
	}
}
`, cfgKeyGoModVersion)))
		xt.OK(t, err)

		m, err := NewGoModVersion(cfg)
		xt.OK(t, err)
		_, err = m.Message()
		xt.KO(t, err)
		xt.Eq(t, "opening go.mod (open my/go.mod: no such file or directory)", err.Error())
	})
}
