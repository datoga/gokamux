package gokamux

type Module interface {
	Execute(ctx Context, message interface{}, params []string)
}
