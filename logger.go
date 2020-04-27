package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type logger struct {
	producer   *kafka.Producer
	logChannel string
}

// Logger struct
type Logger interface {
	WriteLog(interface{}, string) error
}

func (l *logger) produceLog(log interface{}) error {

	// Get schema from channel
	schema, errGetSchema := getLogSchema()
	if errGetSchema != nil {
		return errGetSchema
	}

	// Get codec from schema
	codec, errCreateCodec := getCodec(*schema)
	if errCreateCodec != nil {
		return errCreateCodec
	}

	// Encoding message
	logEncoded, errParseAvro := codec.BinaryFromNative(nil, log)
	if errParseAvro != nil {
		return errors.New("[LOG ENCODING] " + errParseAvro.Error())
	}

	// // We append 5 extra bytes to the beginning of the messages
	// // in order to confluent avro serializer be able to deserrialize
	// // messages in other languages as python.
	// // see: https://docs.confluent.io/current/schema-registry/serializer-formatter.html#wire-format
	// // in section Wire Format.
	// magicByte := []byte{0}
	// schemaBytes := make([]byte, 4)
	// binary.BigEndian.PutUint32(schemaBytes, uint32(schema.ID))
	// wireBytes := append(magicByte, schemaBytes...)
	// logEncoded = append(wireBytes, logEncoded...)

	// Producing message
	l.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &l.logChannel,
			Partition: kafka.PartitionAny,
		},
		Value: logEncoded,
	}

	return nil
}

// WriteLog writes a log on kafka
func (l *logger) WriteLog(message interface{}, action string) error {

	go deliveryReport(l.producer)

	// TODO Spec log's fields
	// Producing log
	logMessage := map[string]interface{}{
		"action":    action,
		"size":      len(fmt.Sprintf("%v", message)),
		"timestamp": time.Now().String(),
	}

	// Kafka insertion
	if errProduceLog := l.produceLog(logMessage); errProduceLog != nil {
		return errProduceLog
	}

	l.producer.Flush(15 * 1000)

	return nil
}

// NewLogger retuns a Logger
func NewLogger() (Logger, error) {

	var l logger

	// Kafka Log Producer Client
	newLogProducer, errKfkLogProducer := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": envs.KafkaBootstrapServers,
	})
	if errKfkLogProducer != nil {
		return nil, errors.New("[LOG PRODUCER] " + errKfkLogProducer.Error())
	}
	l.producer = newLogProducer
	l.logChannel = envs.ChimeraLogChannel

	return &l, nil
}
