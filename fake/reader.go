package fake

import "gitlab.inspr.com/chimera/client-golang"

type reader struct {
	chans []string
}

func (r *reader) ReadMessage() (*string, interface{}, error) {
	messagec := make(chan struct {
		ch *string
		ms interface{}
	})
	for _, c := range r.chans {
		test := c
		go func() {
			for {
				if receivingChannel, ok := channels[test]; ok {
					messagec <- struct {
						ch *string
						ms interface{}
					}{
						&test,
						<-receivingChannel,
					}
				}
			}
		}()
	}
	result := <-messagec
	return result.ch, result.ms, nil
}

func (r *reader) Commit() error {
	return nil
}
func (r *reader) Close() {

}

// NewSimpleReader creates a simple reader that just alters objects
func NewSimpleReader(chans []string) client.Reader {
	return &reader{
		chans: chans,
	}
}
