package rendering

import (
	"bytes"
	"html/template"
)

type HTMLRenderer struct {
	CustomFuncs template.FuncMap
	Layout      []string
}

func (r HTMLRenderer) Render(ctx RenderContext) (output string, err error) {
	files := []string{}
	files = append(files, r.Layout...)
	files = append(files, ctx.Template)

	tmpl := template.New("root").Funcs(r.CustomFuncs)
	tmpl, err = tmpl.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx.Data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
