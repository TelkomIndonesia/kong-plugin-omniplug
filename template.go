package main

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Template struct {
	*template.Template

	text []byte
}

func NewTemplate(s string) (Template, error) {
	t := Template{}
	return t, t.UnmarshalText([]byte(s))
}

func (t Template) ToString(data interface{}) (s string, err error) {
	if t.Empty() {
		return string(t.text), nil
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}
	return buf.String(), nil
}

func (t Template) Empty() bool {
	return t.Template == nil
}

func (t *Template) UnmarshalText(text []byte) (err error) {
	t.Template, err = template.New("template").Funcs(sprig.FuncMap()).Parse(string(text))
	t.text = text
	return
}

func (t Template) MarshalText() (text []byte, err error) {
	return t.text, nil
}
