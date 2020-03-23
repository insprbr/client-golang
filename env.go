package client

import (
	"os"
)

// Env vars struct
type Env struct {
	chimeraNodeID         string
	chimeraNamespace      string
	chimeraRegistryURL    string
	chimeraInputChannels  string
	chimeraOutputChannels string
	chimeraLogChannel     string

	kafkaBootstrapServers string
	kafkaAutoOffsetReset  string
	kafkaEnableAutoCommit string
}

func getEnvVars() *Env {
	return &Env{
		chimeraNodeID:         getEnv("CHIMERA_NODE_ID"),
		chimeraNamespace:      getEnv("CHIMERA_NAMESPACE"),
		chimeraLogChannel:     getEnv("CHIMERA_LOG_CHANNEL"),
		chimeraRegistryURL:    getEnv("CHIMERA_REGISTRY_URL"),
		chimeraInputChannels:  getEnv("CHIMERA_INPUT_CHANNELS"),
		chimeraOutputChannels: getEnv("CHIMERA_OUTPUT_CHANNELS"),
		kafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
		kafkaAutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET"),
	}
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found.")
}

var envs *Env

func init() {
	envs = getEnvVars()
}
