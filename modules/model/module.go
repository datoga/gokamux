package model

type Module interface {
	Init(params ...string) error
	Process(ctx Context, msg string)
}
