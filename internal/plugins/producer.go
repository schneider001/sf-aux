package plugins

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"os"
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
	config.Producer.Return.Errors = true

	config.Producer.Retry.Max = 1
	config.Producer.Retry.Backoff = 3 * time.Second
	config.Producer.Timeout = 10 * time.Second
	config.Producer.RequiredAcks = -1
	config.Net.MaxOpenRequests = 1
	config.Producer.Idempotent = true

	if os.Getenv("KAFKA_NET_SASL_ENABLE") == "true" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = os.Getenv("KAFKA_NET_SASL_USER")
		config.Net.SASL.Password = os.Getenv("KAFKA_NET_SASL_PASSWORD")
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Handshake = true
	}

	if os.Getenv("KAFKA_NET_SECURITY_TLS_ENABLED") == "true" {
		config.Net.TLS.Enable = true
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}

		caFile := os.Getenv("KAFKA_NET_SECURITY_ROOT_CA_FILE")
		if caFile != "" {
			caCert, err := os.ReadFile(caFile)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read CA file")
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}

		config.Net.TLS.Config = tlsConfig
	}

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

		message, err := ev.MarshalData()
		if err != nil {
			log.Println("Marshal error", err)
			continue
		}

		// Test durability
		/*
			if err := appendToFile("producer_msgs", message); err != nil {
				log.Println("Error writing to file:", err)
			}
		*/
		_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
			Headers: []sarama.RecordHeader{
				{Key: []byte("exporter"), Value: []byte(ev.Header.Exporter)},
			},
			Topic: p.main_topic,
			Value: sarama.ByteEncoder(message),
			Key:   sarama.ByteEncoder(ev.Header.Exporter),
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

//Function to test durability
/*
func appendToFile(filename string, message []byte) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "opening file for appending")
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", message))
	return err
}
*/
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
