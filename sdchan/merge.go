package sdchan

func Merge(chans []interface{}) <-chan interface{} {
	if len(chans) == 0 {
		return nil
	}
	rc := make(chan interface{})
	go func() {
		defer close(rc)
		chans1 := removeChan(chans, -1) // Copy
		for {
			chosen, v, recvOK := ReceiveSelect(chans1)
			if !recvOK {
				chans1 = removeChan(chans1, chosen)
				if len(chans1) == 0 {
					break
				}
			} else {
				rc <- v
			}
		}
	}()
	return rc
}

func removeChan(chans []interface{}, index int) []interface{} {
	if len(chans) == 0 {
		return nil
	}
	r := make([]interface{}, 0)
	for i, c := range chans {
		if i != index {
			r = append(r, c)
		}
	}
	return r
}
