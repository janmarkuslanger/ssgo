package task

type TaskContext struct {
	OutputDir string
}

type Task interface {
	Run(ctx TaskContext) error
	IsCritical() bool
}
