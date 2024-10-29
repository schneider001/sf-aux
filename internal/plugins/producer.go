package plugins

import (
	"encoding/binary"
	"log"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"sf-aux/internal/kafka"
	"sf-aux/internal/models"
)

type KafkaProducerPlugin struct {
	topic    string
	producer sarama.SyncProducer
}

func MakeKafkaProducer(cfgPrefix string) (*KafkaProducerPlugin, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producerConfig, err := kafka.NewProducerConfig(cfgPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "getting producer config")
	}

	producer, err := sarama.NewSyncProducer(producerConfig.Addresses, config)
	if err != nil {
		return nil, errors.Wrap(err, "creating producer")
	}

	log.Println("kafka KafkaProducerPlugin")
	alerter := &KafkaProducerPlugin{
		topic:    producerConfig.Topic,
		producer: producer,
	}

	return alerter, nil
}

func (p *KafkaProducerPlugin) Handle(events <-chan models.EventWithContext) error {
	for {
		ev, ok := <-events
		if !ok {
			log.Println("Channel 'records' closed")
			return nil
		}

		version := make([]byte, 8)
		binary.LittleEndian.PutUint64(version, uint64(ev.Header.Version))

		message, err := ev.MarshalData()
		if err != nil {
			log.Println("Marshal error", err)
			continue
		}

		_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
			Headers: []sarama.RecordHeader{
				{Key: []byte("filename"), Value: []byte(ev.Header.Filename)},
				{Key: []byte("exporter"), Value: []byte(ev.Header.Exporter)},
				{Key: []byte("version"), Value: version},
				{Key: []byte("ip"), Value: []byte(ev.Header.Ip)},
			},
			Topic: p.topic,
			Value: sarama.ByteEncoder(message),
		})
		if err != nil {
			log.Println("Send error: ", err)
			continue
		}

	}
}
