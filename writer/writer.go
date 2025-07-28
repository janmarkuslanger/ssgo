package writer

type Writer interface {
	Write(path string, content string) error
}
