package planner

import (
	"bytes"
	"text/template"
)

type Template struct {
	tmpl *template.Template
}

func NewTemplate(text string) (*Template, error) {
	tmpl, err := template.New("rule").Parse(text)
	if err != nil {
		return nil, err
	}
	return &Template{tmpl: tmpl}, nil
}

func (t *Template) Execute(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
