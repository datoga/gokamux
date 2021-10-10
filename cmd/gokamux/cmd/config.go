package cmd

import (
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/datoga/saramaconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Brokers      []string
	InputTopics  []string
	OutputTopics []string
	Operations   []Operation
}

type Operation struct {
	Module string
	Params []string
}

const defaultBroker = "localhost:9092"
const defaultInputTopic = "test"

func loadGokaMuxConfig(v *viper.Viper) (*Config, error) {
	brokers := v.GetStringSlice("brokers")

	if len(brokers) == 0 {
		fmt.Println("No brokers configured, selecting the default broker", defaultBroker)

		brokers = []string{defaultBroker}
	}

	inputTopics := v.GetStringSlice("stream.input")

	if len(inputTopics) == 0 {
		fmt.Println("No input topics configured, selecting the default input topic", defaultInputTopic)

		inputTopics = []string{defaultInputTopic}
	}

	outputTopics := v.GetStringSlice("stream.output")

	if len(outputTopics) == 0 {
		fmt.Println("No output topics configured")
		outputTopics = []string{}
	}

	operationNames := v.GetStringSlice("stream.operations")

	if operationNames == nil || len(operationNames) == 0 {
		fmt.Println("No operations defined on the stream")
	}

	var operationsDefined []Operation

	for _, operation := range operationNames {
		vOperation := v.Sub(operation)

		if vOperation == nil {
			return nil, fmt.Errorf("operation %s defined on the list but not in the module", operation)
		}

		module := vOperation.GetString("module")

		if module == "" {
			return nil, fmt.Errorf("no module name defined for operation %s", operation)
		}

		params, err := parseParams(vOperation.Get("params"))

		if err != nil {
			return nil, fmt.Errorf("failed parsing params with error %v", err)
		}

		operationsDefined = append(
			operationsDefined,
			Operation{
				Module: module,
				Params: params,
			},
		)
	}

	return &Config{
		Brokers:      brokers,
		InputTopics:  inputTopics,
		OutputTopics: outputTopics,
		Operations:   operationsDefined,
	}, nil
}

func parseParams(rawParams interface{}) ([]string, error) {
	var params []string

	if rawParams == nil {
		return params, nil
	}

	switch typeParams := rawParams.(type) {
	case []interface{}:
		for _, v := range typeParams {
			s, ok := v.(string)

			if !ok {
				return nil, fmt.Errorf("could not parse param array %v, wrong value %v", typeParams, v)
			}

			params = append(params, s)
		}
	case string:
		params = append(params, typeParams)
	case int:
		params = append(params, strconv.Itoa(typeParams))
	default:
		return nil, fmt.Errorf("not able to parse params value %v with type %T", typeParams, typeParams)
	}

	return params, nil
}

func loadSaramaConfig(v *viper.Viper) (*sarama.Config, error) {
	saramaCfg, err := saramaconfig.NewFromViper(v)

	if err != nil {
		return nil, fmt.Errorf("failed configuring Sarama with error %v", err)
	}

	return saramaCfg, nil
}
