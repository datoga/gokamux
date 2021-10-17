package gokamux

import (
	"context"
	"fmt"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type Step struct {
	ID     string
	Module string
	Params []string
}

type Muxer struct {
	brokers       []string
	inputs        []string
	outputs       []string
	steps         []Step
	consumerGroup string
}

func NewMuxer() *Muxer {
	return &Muxer{}
}

func (m *Muxer) ConsumerGroup(consumerGroup string) *Muxer {
	m.consumerGroup = consumerGroup

	return m
}

func (m *Muxer) Brokers(broker ...string) *Muxer {
	m.brokers = append(m.brokers, broker...)

	return m
}

func (m *Muxer) Input(topic ...string) *Muxer {
	m.inputs = append(m.inputs, topic...)

	return m
}

func (m *Muxer) Output(topic ...string) *Muxer {
	m.outputs = append(m.outputs, topic...)

	return m
}

func (m *Muxer) Step(steps ...Step) {
	m.steps = append(m.steps, steps...)
}

func (m Muxer) Run(ctx context.Context) error {
	p, err := newPipeline(m.steps...).Compile()

	if err != nil {
		return err
	}

	cb := func(ctx goka.Context, msg interface{}) {
		message, ok := msg.(string)

		if !ok {
			ctx.Fail(fmt.Errorf("failed converting message into string"))
			return
		}

		result := p.Run(ctx, &message)

		if result.Error != nil {
			ctx.Fail(result.Error)
			return
		}

		if result.Discard {
			return
		}

		for _, o := range m.outputs {
			ctx.Emit(goka.Stream(o), ctx.Key(), message)
		}
	}

	var inputs []goka.Edge

	for _, input := range m.inputs {
		inputs = append(inputs, goka.Input(goka.Stream(input), new(codec.String), cb))
	}

	var outputs []goka.Edge

	for _, output := range m.outputs {
		outputs = append(outputs, goka.Output(goka.Stream(output), new(codec.String)))
	}

	var edges []goka.Edge

	edges = append(edges, inputs...)
	edges = append(edges, outputs...)

	var group = m.consumerGroup

	if group == "" {
		group = "test"
	}

	g := goka.DefineGroup(goka.Group(group), edges...)

	processor, err := goka.NewProcessor(m.brokers, g)

	if err != nil {
		return fmt.Errorf("failed initialising processor with error %v", err)
	}

	return processor.Run(ctx)
}
