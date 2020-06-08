package client

import (
	"errors"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// NewSingleChannelReader returns a new reader
func NewSingleChannelReader(channel string) (Reader, error) {
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
	channelsToConsume := func() []string {
		ret := []string{}
		for _, s := range strings.Split(listOfChannels, ";") {
			ret = append(ret, toTopic(s))
		}
		return ret
	}()

	if !validChannel(toTopic(channel), channelsToConsume) {
		return nil, errors.New("invalid channel")
	}

	r.consumer.Subscribe(toTopic(channel), nil)
	return &r, nil
}
