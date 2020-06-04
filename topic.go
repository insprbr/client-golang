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
		ch.name = topic[len("chimera-"+ch.namespace+"-"):]
	} else {
		ch.name = topic[len("chimera-"+ch.prefix+"-"+ch.namespace+"-"):]
	}
	return ch
}

func toTopic(ch string) string {
	var topic string

	if envs.ChimeraEnvironment == "" {
		topic = fmt.Sprintf("chimera-%s-%s", envs.ChimeraNamespace, envs.ChimeraNodeID)
	} else {
		topic = fmt.Sprintf(
			"chimera-%s-%s-%s",
			envs.ChimeraEnvironment,
			envs.ChimeraNamespace,
			envs.ChimeraNodeID,
		)
	}

	return topic
}
