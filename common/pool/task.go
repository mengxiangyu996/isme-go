package pool

type Task interface {
	Execute() error
}