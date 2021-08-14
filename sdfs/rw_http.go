package sdfs

import (
	"io/ioutil"
	"net/http"

	"github.com/gaorx/stardust3/sderr"
)

func HttpReadBytes(hfs http.FileSystem, name string) ([]byte, error) {
	if hfs == nil {
		return nil, sderr.New("nil hfs")
	}
	f, err := hfs.Open(name)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() { _ = f.Close() }()
	r, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func HttpReadText(hfs http.FileSystem, name string) (string, error) {
	b, err := HttpReadBytes(hfs, name)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return string(b), nil
}

func HttpReadTextDef(hfs http.FileSystem, name, def string) string {
	s, err := HttpReadText(hfs, name)
	if err != nil {
		return def
	}
	return s
}
