package client

import (
	"errors"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type reader struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
	logger      Logger
}

// Reader reads messages from Chimera
type Reader interface {
	Commit() error
	ReadMessage() (*string, interface{}, error)
	Close()
}

// Commit commits the last message read by Reader
func (r *reader) Commit() error {
	_, errCommit := r.consumer.CommitMessage(r.lastMessage)
	if errCommit != nil {
		return errors.New("[READER_COMMIT] " + errCommit.Error())
	}
	return nil
}

//
// ReadMessage reads message by message
// Returns channel the message belongs to, the message and an error if any occured.
func (r *reader) ReadMessage() (*string, interface{}, error) {

	for {
		select {
		default:
			// TODO see the other way to read message from kafka
			ev := r.consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:

				channel := *e.TopicPartition.Topic
				log.Println("reading message from channel " + channel)
				// Decoding Message
				message, errDecode := decode(e.Value, fromTopic(channel))
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
				log.Println(e)

			default:
				continue
			}
		}
	}
}

// Close reader consumer
func (r *reader) Close() {
	r.consumer.Close()
}

// NewReader returns a new reader
func NewReader() (Reader, error) {
	envs = GetEnvVars()
	// Reader consumer Client
	var r reader

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  envs.KafkaBootstrapServers,
		"group.id":           envs.ChimeraNodeID,
		"auto.offset.reset":  envs.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, errors.New("[NEW READER] " + errKafkaConsumer.Error())
	}
	r.consumer = newConsumer

	// Logger
	newLogger, errNewLogger := NewLogger()
	if errNewLogger != nil {
		return nil, errors.New("[NEW READER] " + errNewLogger.Error())
	}
	r.logger = newLogger

	// Get channels to consume and subscribe
	listOfChannels := envs.ChimeraInputChannels
	if listOfChannels == "" {
		return nil, errors.New("[ENV VAR] KAFKA_INPUT_CHANNELS not specified. ")
	}
	channelList := strings.Split(listOfChannels, ";")
	channelList = channelList[:len(channelList)-1]
	channelsToConsume := func() []string {
		ret := []string{}
		for _, s := range channelList {
			topic := toTopic(s)
			log.Println(topic)
			ret = append(ret, topic)
		}
		return ret
	}()
	r.consumer.SubscribeTopics(channelsToConsume, nil)
	tpoics, _ := r.consumer.Subscription()
	log.Println(tpoics)
	return &r, nil
}
