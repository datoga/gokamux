package cmd

import (
	"fmt"

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
	Type   string
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

	operations := v.Sub("stream.operations")

	operationsDefined := []Operation{}

	if operations == nil {
		fmt.Println("No operations defined on the stream")
	}

	return &Config{
		Brokers:      brokers,
		InputTopics:  inputTopics,
		OutputTopics: outputTopics,
		Operations:   operationsDefined,
	}, nil
}

func loadSaramaConfig(v *viper.Viper) (*sarama.Config, error) {
	saramaCfg, err := saramaconfig.NewFromViper(v)

	if err != nil {
		return nil, fmt.Errorf("failed configuring Sarama with error %v", err)
	}

	return saramaCfg, nil
}
