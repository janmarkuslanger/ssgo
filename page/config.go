package page

type PagePayload struct {
	Params map[string]string
	Path   string
}

type Config struct {
	Template string
	Pattern  string
	GetData  func(payload PagePayload) map[string]any
	GetPaths func() map[string]string
}
