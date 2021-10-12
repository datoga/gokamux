package gokamux

type ModuleLoader interface {
	Init(params ...string) Module
}

type Module interface {
	Execute(ctx Context, msg string)
}
