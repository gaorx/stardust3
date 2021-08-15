package sdhttptmpl

import (
	"bytes"
	"text/template"

	"github.com/gaorx/stardust3/sderr"
)

func ExecStr(tmpl string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return buff.String(), nil
}

func MustExecStr(tmpl string, data interface{}) string {
	r, err := ExecStr(tmpl, data)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return r
}

func ExecStrDef(tmpl string, data interface{}, def string) string {
	r, err := ExecStr(tmpl, data)
	if err != nil {
		return def
	}
	return r
}
