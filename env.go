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
	chimeraEnvironment    string

	kafkaBootstrapServers string
	kafkaAutoOffsetReset  string
	kafkaEnableAutoCommit string
}

// GetEnvVars returns a struct with all Chimera envirnoment variables
func GetEnvVars() *Env {
	return &Env{
		chimeraNodeID:         getEnv("CHIMERA_NODE_ID"),
		chimeraNamespace:      getEnv("CHIMERA_NAMESPACE"),
		chimeraLogChannel:     getEnv("CHIMERA_LOG_CHANNEL"),
		chimeraRegistryURL:    getEnv("CHIMERA_REGISTRY_URL"),
		chimeraInputChannels:  getEnv("CHIMERA_INPUT_CHANNELS"),
		chimeraOutputChannels: getEnv("CHIMERA_OUTPUT_CHANNELS"),
		kafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
		kafkaAutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET"),
		chimeraEnvironment:    getEnv("CHIMERA_ENVIRONMENT"),
	}
}
func testEnvVars(
	name string,
	namespace string,
	log string,
	registry string,
	input string,
	output string,
	environment string,
) *Env {
	return &Env{
		chimeraNodeID:         name,
		chimeraNamespace:      namespace,
		chimeraLogChannel:     log,
		chimeraRegistryURL:    registry,
		chimeraInputChannels:  input,
		chimeraOutputChannels: output,
		kafkaBootstrapServers: "kafka",
		kafkaAutoOffsetReset:  "latest",
		chimeraEnvironment:    environment,
	}
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found.")
}

var envs *Env
