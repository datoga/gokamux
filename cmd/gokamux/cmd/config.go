package cmd

import (
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/datoga/gokamux"
	"github.com/datoga/saramaconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Brokers      []string
	InputTopics  []string
	OutputTopics []string
	Steps        []gokamux.Step
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

	stepNames := v.GetStringSlice("stream.steps")

	if stepNames == nil || len(stepNames) == 0 {
		fmt.Println("No steps defined on the stream")
	}

	var steps []gokamux.Step

	for _, stepName := range stepNames {
		vStep := v.Sub(stepName)

		if vStep == nil {
			return nil, fmt.Errorf("step %s defined on the list but not in the module", stepName)
		}

		module := vStep.GetString("module")

		if module == "" {
			return nil, fmt.Errorf("no module name defined for step %s", stepName)
		}

		params, err := parseParams(vStep.Get("params"))

		if err != nil {
			return nil, fmt.Errorf("failed parsing params with error %v", err)
		}

		steps = append(
			steps,
			gokamux.Step{
				ID:     stepName,
				Module: module,
				Params: params,
			},
		)
	}

	return &Config{
		Brokers:      brokers,
		InputTopics:  inputTopics,
		OutputTopics: outputTopics,
		Steps:        steps,
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
