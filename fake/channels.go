package fake

var channels map[string]chan interface{}

func init() {
	channels = map[string]chan interface{}{}
}
