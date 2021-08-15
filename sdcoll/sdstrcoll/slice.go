package sdstrcoll

type Slice []string

func (l Slice) Has(v string) bool {
	for _, v0 := range l {
		if v0 == v {
			return true
		}
	}
	return false
}

func (l Slice) Strings() []string {
	return []string(l)
}

func (l Slice) Copy() Slice {
	if l == nil {
		return nil
	}
	l1 := make(Slice, 0, len(l))
	l1 = append(l1, l...)
	return l1
}
