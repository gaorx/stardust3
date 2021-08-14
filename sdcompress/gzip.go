package sdcompress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/gaorx/stardust3/sderr"
)

func Gzip(data []byte, level int) ([]byte, error) {
	if data == nil {
		return nil, sderr.New("nil data")
	}

	buff := new(bytes.Buffer)
	w, err := gzip.NewWriterLevel(buff, level)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	err = w.Close()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buff.Bytes(), nil
}

func Ungzip(data []byte) ([]byte, error) {
	if data == nil {
		return nil, sderr.New("nil data")
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() { _ = r.Close() }()

	to, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return to, nil
}
