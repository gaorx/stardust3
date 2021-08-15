package sdtime

import (
	"github.com/gaorx/stardust3/sdparse"
)

var (
	Parse     = sdparse.Time
	ParseDef  = sdparse.TimeDef
	MustParse = sdparse.MustTime
)
