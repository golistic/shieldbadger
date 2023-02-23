// Copyright (c) 2023, Geert JM Vanderkelen

package shieldbadger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golistic/xt"
)

func TestNewConfigFromFile(t *testing.T) {
	t.Run("file doesn't exist", func(t *testing.T) {
		_, err := newConfigFromFile("something/that/doesnot/exist.json")
		xt.KO(t, err)
	})

	t.Run("incorrect JSON", func(t *testing.T) {
		data := bytes.NewBufferString("{not:json}")

		_, err := newConfigFromReader(data)
		xt.KO(t, err)
		xt.Eq(t, "invalid character 'n' looking for beginning of object key string", err.Error())
	})

	t.Run("from file", func(t *testing.T) {
		cfgFile := filepath.Join("_testdata", "akdfieifksk.json")
		fp, err := os.Create(cfgFile)
		xt.OK(t, err)
		defer func() { _ = os.Remove(cfgFile) }()

		_, err = fp.WriteString(`
{
	"style": "plastic",
	"messengers": {
		"go.mod.version": {
			"path": "my/go.mod"
		}
	}
}
`)
		xt.OK(t, err)
		xt.OK(t, fp.Close())

		cfg, err := newConfigFromFile(cfgFile)
		xt.OK(t, err)
		xt.Eq(t, "plastic", cfg.Style)
		xt.Eq(t, "json.RawMessage", fmt.Sprintf("%T", cfg.Messengers["go.mod.version"]))
	})
}
