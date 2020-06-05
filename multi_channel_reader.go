package client

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

// MultiChannelReader is a reader that reads from multiple channels
type MultiChannelReader interface {
	ReadMessage(c string) (interface{}, error)
	Commit(c string) error
	CommitAll() error
}

var mcr *multiChannelReader

type multiChannelReader struct {
	readers map[string]Reader
}

// NewMultiChannelReader returns a multi channel reader.
func NewMultiChannelReader() (MultiChannelReader, error) {
	if mcr == nil {
		log.Println("instantiating new multi channel reader")
		mcr = &multiChannelReader{
			map[string]Reader{},
		}
		envs = GetEnvVars()
		log.Println("environment:")
		prettyEnviron, _ := json.MarshalIndent(envs, "", "\t")
		log.Println(prettyEnviron)
		// Get channels to consume and subscribe
		listOfChannels := envs.ChimeraInputChannels
		if listOfChannels == "" {
			return nil, errors.New("[ENV VAR] KAFKA_INPUT_CHANNELS not specified. ")
		}

		channelsToConsume := strings.Split(listOfChannels, ";")
		channelsToConsume = channelsToConsume[:len(channelsToConsume)-1]
		log.Println("input channels:")
		log.Println(channelsToConsume)
		for _, c := range channelsToConsume {
			var errReaderInstanciation error
			mcr.readers[c], errReaderInstanciation = NewSingleChannelReader(c)
			if errReaderInstanciation != nil {
				return nil, errReaderInstanciation
			}
		}
	}
	return mcr, nil
}

func (r *multiChannelReader) ReadMessage(c string) (interface{}, error) {
	singleReader, ok := r.readers[c]
	if !ok {
		return nil, errors.New("[CLIENT] invalid channel in multi channel reader")
	}
	_, message, err := singleReader.ReadMessage()
	return message, err
}

func (r *multiChannelReader) Commit(c string) error {
	singleReader, ok := r.readers[c]
	if !ok {
		return errors.New("[CLIENT] invalid channel in multi channel reader")
	}
	err := singleReader.Commit()
	return err
}

func (r *multiChannelReader) CommitAll() error {
	for _, c := range r.readers {
		err := c.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}
