package sdiploc

import (
	"bytes"
	_ "embed"

	"github.com/gaorx/stardust3/sderr"
)

//go:embed loc.csv
var locBytes []byte

//go:embed rec.csv
var recordBytes []byte

func init() {
	finder, err := Load(bytes.NewReader(recordBytes), bytes.NewReader(locBytes))
	if err != nil {
		panic(sderr.WithStack(err))
	}
	DefaultFinder = finder
}
