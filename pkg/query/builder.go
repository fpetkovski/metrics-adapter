package query

import (
	"bytes"
	"text/template"
)

type Builder struct {
}

func NewQueryBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildQuery(tpl string, data interface{}) string {
	t := template.New("")
	t, err := t.Parse(tpl)
	if err != nil {
		return ""
	}

	var result bytes.Buffer
	if err := t.Execute(&result, data); err != nil {
		return ""
	}

	return result.String()
}
