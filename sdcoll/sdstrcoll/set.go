package sdstrcoll

import (
	"encoding/json"
	"fmt"
)

type Set struct {
	data map[string]int
}

func (set Set) String() string {
	return fmt.Sprintf("%s", set.Elems())
}

func (set Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Elems())
}

func (set *Set) UnmarshalJSON(data []byte) error {
	if set == nil {
		return nil
	}
	var l []string
	err := json.Unmarshal(data, &l)
	if err != nil {
		return err
	}
	set.ensure()
	for _, k := range l {
		set.data[k] = 0
	}
	return nil
}

func SetOf(elems ...string) Set {
	set := Set{}
	set.Add(elems...)
	return set
}

func (set Set) Elems() []string {
	elems := make([]string, 0, len(set.data))
	for elem := range set.data {
		elems = append(elems, elem)
	}
	return elems
}

func (set *Set) ensure() {
	if set.data == nil {
		set.data = make(map[string]int)
	}
}

func (set *Set) Add(elems ...string) {
	set.ensure()
	for _, elem := range elems {
		set.data[elem] = 0
	}
}

func (set *Set) Del(elems ...string) {
	if len(set.data) == 0 || len(elems) == 0 {
		return
	}
	for _, elem := range elems {
		delete(set.data, elem)
	}
}

func (set *Set) Union(others ...Set) {
	if len(others) == 0 {
		return
	}
	set.ensure()
	for _, other := range others {
		for elem := range other.data {
			set.data[elem] = 0
		}
	}
}

func (set Set) Size() int {
	return len(set.data)
}

func (set Set) Has(elem string) bool {
	_, ok := set.data[elem]
	return ok
}

func (set Set) HasOne(elems ...string) bool {
	for _, elem := range elems {
		_, ok := set.data[elem]
		if ok {
			return true
		}
	}
	return false
}

func (set Set) HasAll(elems ...string) bool {
	for _, elem := range elems {
		_, ok := set.data[elem]
		if !ok {
			return false
		}
	}
	return true
}

func ReduceSet(accumulator, each Set) Set {
	accumulator.Add(each.Elems()...)
	return accumulator
}
