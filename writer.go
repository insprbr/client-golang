package client

import (
	"errors"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type writer struct {
	producer *kafka.Producer
	logger   Logger
}

// Writer writes messages to Chimera channels
type Writer interface {
	WriteMessage(interface{}, string) error
}

func (w *writer) produceMessage(message interface{}, channel string) error {

	messageEncoded, errorEncode := encode(message, channel)
	if errorEncode != nil {
		return errorEncode
	}

	channel = envs.chimeraNamespace + "_" + channel

	// Producing message
	w.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &channel,
			Partition: kafka.PartitionAny,
		},
		Value: messageEncoded,
	}

	return nil
}

// Messenger report
func deliveryReport(producer *kafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
			return
		default:
			log.Println(ev)
		}
	}
}

func validChannel(channel string, outputChannels []string) bool {
	for _, ch := range outputChannels {
		if ch == channel {
			return true
		}
	}
	return false
}

// WriteMessage writes a message on kafka
func (w *writer) WriteMessage(message interface{}, channel string) error {
	// TODO Check if it's a valid channel
	outputChannels := strings.Split(envs.chimeraOutputChannels, ";")
	if !validChannel(channel, outputChannels) {
		return errors.New("[OUTPUT CHANNEL] Invalid channel. ")
	}

	go deliveryReport(w.producer)

	// Kafka insertion
	if errProduceMessage := w.produceMessage(message, channel); errProduceMessage != nil {
		return errProduceMessage
	}

	if errProduceLog := w.logger.WriteLog(message, "Write"); errProduceLog != nil {
		return errProduceLog
	}

	w.producer.Flush(15 * 1000)

	return nil
}

// Close wtriter producer
func (w *writer) Close() {
	w.producer.Close()
}

// NewWriter returns a new writer
func NewWriter() (Writer, error) {
	envs = getEnvVars()
	var w writer

	// Kafka Producer Client
	newProducer, errKfkProducer := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": envs.kafkaBootstrapServers,
	})

	if errKfkProducer != nil {
		return nil, errors.New("[NEW WRITER] " + errKfkProducer.Error())
	}
	w.producer = newProducer

	// Logger
	newLogger, errNewLogger := NewLogger()
	if errNewLogger != nil {
		return nil, errors.New("[NEW WRITER] " + errNewLogger.Error())
	}
	w.logger = newLogger

	return &w, nil
}
