package sdrand

import (
	"math/rand"
	"reflect"

	"github.com/gaorx/stardust3/sderr"
)

func Shuffle(arr interface{}) {
	typ := reflect.TypeOf(arr)
	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		panic(sderr.WithStack(sderr.ErrIllegalType))
	}

	arrVar := reflect.ValueOf(arr)
	n := arrVar.Len()
	if n <= 1 {
		return
	}
	clone := make([]interface{}, n)
	for i := 0; i < n; i++ {
		clone[i] = arrVar.Index(i).Interface()
	}
	perms := rand.Perm(n)
	for i := 0; i < n; i++ {
		newIndex := perms[i]
		v1 := clone[newIndex]
		arrVar.Index(i).Set(reflect.ValueOf(v1))
	}
	return
}

// TODO:ShuffleStrs
// TODO:ShuffleInts
// TODO:ShuffleInt64s
