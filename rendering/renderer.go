package rendering

type RenderContext struct {
	Data     map[string]any
	Template string
}

type Renderer interface {
	Render(ctx RenderContext) (output string, err error)
}
