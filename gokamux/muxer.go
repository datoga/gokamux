package gokamux

import (
	"context"
	"fmt"
)

type Operation struct {
	Name   string
	Module string
	Params []string
}

type Muxer struct {
	brokers    []string
	inputs     []string
	outputs    []string
	operations []Operation
}

func New() *Muxer {
	return &Muxer{}
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

func (m *Muxer) Operation(operations ...Operation) {
	m.operations = append(m.operations, operations...)
}

func (m Muxer) Run(ctx context.Context) error {
	_, err := m.build()

	if err != nil {
		return fmt.Errorf("failed building module pipeline with error %v", err)
	}

	//p.Run()

	return nil
}

func (m Muxer) build() (*pipeline, error) {
	p, err := newPipeline("error")

	return p, err
}
