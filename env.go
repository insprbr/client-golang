package client

import (
	"os"
)

// Env vars struct
type Env struct {
	ChimeraNodeID         string
	ChimeraNamespace      string
	ChimeraRegistryURL    string
	ChimeraInputChannels  string
	ChimeraOutputChannels string
	ChimeraLogChannel     string
	ChimeraEnvironment    string

	KafkaBootstrapServers string
	KafkaAutoOffsetReset  string
}

// GetEnvVars returns a struct with all Chimera envirnoment variables
func GetEnvVars() *Env {
	return &Env{
		ChimeraNodeID:         getEnv("CHIMERA_NODE_ID"),
		ChimeraNamespace:      getEnv("CHIMERA_NAMESPACE"),
		ChimeraLogChannel:     getEnv("CHIMERA_LOG_CHANNEL"),
		ChimeraRegistryURL:    getEnv("CHIMERA_REGISTRY_URL"),
		ChimeraInputChannels:  getEnv("CHIMERA_INPUT_CHANNELS"),
		ChimeraOutputChannels: getEnv("CHIMERA_OUTPUT_CHANNELS"),
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaAutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET"),
		ChimeraEnvironment:    getEnv("CHIMERA_ENVIRONMENT"),
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
		ChimeraNodeID:         name,
		ChimeraNamespace:      namespace,
		ChimeraLogChannel:     log,
		ChimeraRegistryURL:    registry,
		ChimeraInputChannels:  input,
		ChimeraOutputChannels: output,
		KafkaBootstrapServers: "kafka",
		KafkaAutoOffsetReset:  "latest",
		ChimeraEnvironment:    environment,
	}
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found.")
}

var envs *Env
