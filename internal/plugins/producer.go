package plugins

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"sf-aux/internal/kafka"
	"sf-aux/internal/models"
)

type KafkaProducerPlugin struct {
	main_topic        string
	keepalive_topic   string
	producer          sarama.SyncProducer
	lastSentTimestamp time.Time
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
		main_topic:        producerConfig.MainTopic,
		keepalive_topic:   producerConfig.KeepAliveTopic,
		producer:          producer,
		lastSentTimestamp: time.Time{},
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
				{Key: []byte("exporter"), Value: []byte(ev.Header.Exporter)},
			},
			Topic: p.main_topic,
			Value: sarama.ByteEncoder(message),
		})
		if err != nil {
			log.Println("Send error: ", err)
			continue
		}

		now := time.Now()
		if now.Sub(p.lastSentTimestamp) > time.Hour {
			err = p.sendKeepAlive(ev.Header.Exporter)
			if err != nil {
				log.Println("Error sending keep-alive message:", err)
			}
			p.lastSentTimestamp = now
		}
	}
}

func (p *KafkaProducerPlugin) sendKeepAlive(exporter string) error {
	isoTimestamp := time.Now().Format(time.RFC3339)

	keepAliveData := map[string]string{
		"exporter":  exporter,
		"timestamp": isoTimestamp,
	}

	message, err := json.Marshal(keepAliveData)
	if err != nil {
		return errors.Wrap(err, "marshal keep-alive data")
	}

	keepAliveMessage := &sarama.ProducerMessage{
		Topic: p.keepalive_topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err = p.producer.SendMessage(keepAliveMessage)
	return err
}
