package kafka

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type KafkaProducerConfig struct {
	Addresses []string `required:"true"`
	Topic     string   `required:"true"`
}

func NewProducerConfig(envPrefix string) (*KafkaProducerConfig, error) {
	var config KafkaProducerConfig

	if err := envconfig.Process(envPrefix, &config); err != nil {
		return nil, errors.Wrap(err, "processing env")
	}

	return &config, nil
}
