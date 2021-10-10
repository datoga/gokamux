package gokamux

import "github.com/lovoo/goka"

type Context interface {
	Discard()
	OverrideMessage(interface{})
	GokaContext() goka.Context
}

type cbContext struct {
	GokaCtx          goka.Context
	Discarded        bool
	OverridedMessage interface{}
}

func (ct *cbContext) Discard() {
	ct.Discarded = true
}

func (ct *cbContext) OverrideMessage(message interface{}) {
	ct.OverridedMessage = message
}

func (ct cbContext) GokaContext() goka.Context {
	return ct.GokaCtx
}
