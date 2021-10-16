package modules

type ModuleLoader interface {
	Init(params ...string) (Module, error)
}

type Module interface {
	Process(ctx Context, msg string)
}
