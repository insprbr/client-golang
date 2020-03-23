package client

import (
	"errors"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Reader struct
type Reader struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
	logger      *Logger
}

// Commit commits the last message read by Reader
func (r *Reader) Commit() error {
	_, errCommit := r.consumer.CommitMessage(r.lastMessage)
	if errCommit != nil {
		return errors.New("[READER_COMMIT] " + errCommit.Error())
	}
	return nil
}

//
// ReadMessage reads message by message
// Returns channel the message belongs to, the message and an error if any occured.
func (r *Reader) ReadMessage() (*string, interface{}, error) {

	for {
		select {
		default:
			// TODO see the other way to read message from kafka
			ev := r.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:

				channel := *e.TopicPartition.Topic

				// Decoding Message
				message, errDecode := decode(e.Value, channel)
				if errDecode != nil {
					return nil, nil, errDecode
				}

				// Write log
				errWritelog := r.logger.WriteLog(message, "Read")
				if errWritelog != nil {
					return nil, nil, errWritelog
				}

				r.lastMessage = e

				return &channel, message, nil

			case kafka.Error:
				if e.Code() == kafka.ErrAllBrokersDown {
					return nil, nil, errors.New("[FATAL_ERROR]\n===== All brokers are down! =====\n" + e.Error())
				}
			default:
				continue
			}
		}
	}
}

// Close Reader consumer
func (r *Reader) Close() {
	r.consumer.Close()
}

// NewReader returns a new Reader
func NewReader() (*Reader, error) {

	// Reader consumer Client
	var reader Reader

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  envs.kafkaBootstrapServers,
		"group.id":           envs.chimeraNodeID,
		"auto.offset.reset":  envs.kafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, errors.New("[NEW READER] " + errKafkaConsumer.Error())
	}
	reader.consumer = newConsumer

	// Logger
	newLogger, errNewLogger := NewLogger()
	if errNewLogger != nil {
		return nil, errors.New("[NEW READER] " + errNewLogger.Error())
	}
	reader.logger = newLogger

	// Get channels to consume and subscribe
	listOfChannels := envs.chimeraInputChannels
	if listOfChannels == "" {
		return nil, errors.New("[ENV VAR] KAFKA_INPUT_CHANNELS not specified. ")
	}
	channelsToConsume := func() []string {
		ret := []string{}
		for _, s := range strings.Split(listOfChannels, ";") {
			ret = append(ret, envs.chimeraNamespace+"_"+s)
		}
		return ret
	}()
	reader.consumer.SubscribeTopics(channelsToConsume, nil)

	return &reader, nil
}
