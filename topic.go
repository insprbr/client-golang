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
		namespace: envs.ChimeraNamespace,
		prefix:    envs.ChimeraEnvironment,
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

	if envs.ChimeraEnvironment == "" {
		topic = fmt.Sprintf("chimera_%s_%s", envs.ChimeraNamespace, envs.ChimeraNodeID)
	} else {
		topic = fmt.Sprintf(
			"chimera_%s_%s_%s",
			envs.ChimeraEnvironment,
			envs.ChimeraNamespace,
			envs.ChimeraNodeID,
		)
	}

	return topic
}
