package sdcall

import (
	"github.com/gaorx/stardust3/sdsync"
)

var (
	Lock  = sdsync.Lock
	LockR = sdsync.LockR
	LockW = sdsync.LockW
)
