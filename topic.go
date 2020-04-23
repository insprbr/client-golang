package client

import (
	"fmt"
)

type channel struct {
	namespace string
	name      string
	prefix    string
}

func fromTopic(topic string) channel {
	ch := channel{
		namespace: envs.chimeraNamespace,
		prefix:    envs.chimeraEnvironment,
	}

	if ch.prefix == "" {
		ch.name = topic[len("chimera_"+ch.namespace+"_"):]
	} else {
		ch.name = topic[len("chimera_"+ch.prefix+"_"+ch.namespace+"_"):]
	}
	return ch
}

func toTopic(ch string) string {
	var topic string

	if envs.chimeraEnvironment == "" {
		topic = fmt.Sprintf("chimera-%s-%s", envs.chimeraNamespace, envs.chimeraNodeID)
	} else {
		topic = fmt.Sprintf(
			"chimera-%s-%s-%s",
			envs.chimeraEnvironment,
			envs.chimeraNamespace,
			envs.chimeraNodeID,
		)
	}

	return topic
}
