package sdbytes

func Copy(data []byte) []byte {
	if data == nil {
		return nil
	}
	n := len(data)
	b := make([]byte, n)
	copy(b, data)
	return b
}
