package fake

import "gitlab.inspr.com/chimera/client-golang"

type writer struct {
}

func (w *writer) WriteMessage(message interface{}, channel string) error {
	if _, ok := channels[channel]; !ok {
		channels[channel] = make(chan interface{}, 1000)
	}
	channels[channel] <- message
	return nil
}

// NewSimpleWriter creates a simple writer that just changes object creation
func NewSimpleWriter() client.Writer {
	return &writer{}
}
