package query

import (
	"bytes"
	"fmt"
	"text/template"
)

type Builder struct {
}

func NewQueryBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildQuery(tpl string, data interface{}) (string, error) {
	t := template.New("")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("unable to render query: %w", err)
	}

	var result bytes.Buffer
	if err := t.Execute(&result, data); err != nil {
		return "", fmt.Errorf("unable to render query: %w", err)
	}

	return result.String(), nil
}
