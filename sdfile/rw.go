package sdfile

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/gaorx/stardust3/sderr"
)

func WriteBytes(filename string, data []byte, perm os.FileMode) error {
	err := ioutil.WriteFile(filename, data, perm)
	return sderr.WithStack(err)
}

func WriteText(filename string, text string, perm os.FileMode) error {
	return WriteBytes(filename, []byte(text), perm)
}

func AppendBytes(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return sderr.WithStack(err)
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	_ = f.Close()
	return sderr.WithStack(err)
}

func AppendText(filename string, text string, perm os.FileMode) error {
	return AppendBytes(filename, []byte(text), perm)
}

func ReadBytes(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return data, nil
}

func ReadBytesDef(filename string, def []byte) []byte {
	data, err := ReadBytes(filename)
	if err != nil {
		return def
	}
	return data
}

func ReadText(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return string(data), nil
}

func ReadTextDef(filename, def string) string {
	data, err := ReadText(filename)
	if err != nil {
		return def
	}
	return data
}
