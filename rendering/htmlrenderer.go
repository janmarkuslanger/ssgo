package rendering

import (
	"bytes"
	"html/template"
)

type HTMLRenderer struct{}

func (r HTMLRenderer) Render(ctx RenderContext) (output string, err error) {
	files := []string{}
	files = append(files, ctx.Layout...)
	files = append(files, ctx.Template)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx.Data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
