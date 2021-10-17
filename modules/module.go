package modules

type Module interface {
	Process(ctx Context, msg string)
}

type Configurer interface {
	Configure(params ...string) error
}
