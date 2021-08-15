package sdiploc

import (
	"github.com/gaorx/stardust3/sderr"
)

var (
	DefaultFinder *Finder = nil
)

func FindLoc(ip string) *Loc {
	finder := DefaultFinder
	if finder == nil {
		panic(sderr.New("nil DefaultFinder"))
	}
	return finder.Loc(ip)
}

func FindLocInt(ip uint32) *Loc {
	finder := DefaultFinder
	if finder == nil {
		panic(sderr.New("nil DefaultFinder"))
	}
	return finder.LocInt(ip)
}
