package gokamux

type ModuleResultNotifier interface {
	Discard()
	OverrideMessage(message string)
}

type ModuleFn func(notifier ModuleResultNotifier, message string, params []string)

type Module struct {
	Name string
	fn   ModuleFn
}

type FilterResult bool

var Allow = FilterResult(false)
var Discard = FilterResult(true)

func New(name string, fn ModuleFn) *Module {
	return &Module{
		Name: name,
		fn:   fn,
	}
}

func (m Module) Execute(message *string, params []string) FilterResult {
	mCtx := moduleCtx{}

	m.fn(&mCtx, *message, params)

	if mCtx.IsDiscarded() {
		return Discard
	}

	if overrides, overrideMessage := mCtx.OverridesMessage(); overrides {
		message = &overrideMessage
	}

	return Allow
}

type moduleCtx struct {
	discard         bool
	overrideMessage string
}

func (mCtx *moduleCtx) Discard() {
	mCtx.discard = true
}

func (mCtx *moduleCtx) OverrideMessage(message string) {
	mCtx.overrideMessage = message
}

func (mCtx moduleCtx) IsDiscarded() bool {
	return mCtx.discard
}

func (mCtx moduleCtx) OverridesMessage() (bool, string) {
	return mCtx.overrideMessage != "", mCtx.overrideMessage
}
