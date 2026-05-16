package prompt

import (
	"bytes"
	"context"
	"text/template"
)

type TemplateRenderer interface {
	Render(ctx context.Context, templateText string, data map[string]interface{}) (string, error)
}

type Renderer struct{}

func NewRenderer() Renderer {
	return Renderer{}
}

func (Renderer) Render(ctx context.Context, templateText string, data map[string]interface{}) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return renderTemplate(templateText, data)
}

func renderTemplate(templateText string, data map[string]interface{}) (string, error) {
	tpl, err := template.New("prompt").Option("missingkey=zero").Parse(templateText)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
