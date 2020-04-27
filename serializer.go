package client

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/linkedin/goavro"
)

func getSchema(ch channel) (*string, error) {

	// Getting schema
	resp, errGetSchema := http.Get(envs.ChimeraRegistryURL + "/schema/" + ch.namespace + "/" + ch.name)
	if errGetSchema != nil {
		return nil, errors.New("[KAFKA_PRODUCE_SCHEMA] " + errGetSchema.Error())
	}
	if resp.StatusCode != 200 {
		log.Println(envs.ChimeraRegistryURL + "/schema/" + ch.namespace + "/" + ch.name)
		return nil, errors.New("[KAFKA_PRODUCE_SCHEMA] Schema not registered. ")
	}
	defer resp.Body.Close()

	bBody, errDecodeBody := ioutil.ReadAll(resp.Body)
	if errDecodeBody != nil {
		return nil, errors.New("[PARSE BODY] " + errDecodeBody.Error())
	}
	defer resp.Body.Close()

	sBody := string(bBody)

	return &sBody, nil
}

func getLogSchema() (*string, error) {

	// Getting schema
	resp, errGetSchema := http.Get(envs.ChimeraRegistryURL + "/schema/" + "chimera" + "/" + "log")
	if errGetSchema != nil {
		return nil, errors.New("[KAFKA_PRODUCE_SCHEMA] " + errGetSchema.Error())
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("[KAFKA_PRODUCE_SCHEMA] Schema not registered. ")
	}
	defer resp.Body.Close()

	bBody, errDecodeBody := ioutil.ReadAll(resp.Body)
	if errDecodeBody != nil {
		return nil, errors.New("[PARSE BODY] " + errDecodeBody.Error())
	}
	defer resp.Body.Close()

	sBody := string(bBody)

	return &sBody, nil
}

func getCodec(schema string) (*goavro.Codec, error) {

	// Creating the codec
	codec, errCreateCodec := goavro.NewCodec(schema)
	if errCreateCodec != nil {
		return nil, errors.New("[KAFKA_PRODUCE_CODEC] " + errCreateCodec.Error())
	}

	return codec, nil
}

func decode(messageEncoded []byte, ch channel) (interface{}, error) {

	// Get schema from ch
	schema, errGetSchema := getSchema(ch)
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	// Get codec from schema
	codec, errCreateCodec := getCodec(*schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	// // Decoding message
	// // Removing 5 extra bytes
	// messageEncoded = messageEncoded[5:]

	message, _, errDecoding := codec.NativeFromBinary(messageEncoded)
	if errDecoding != nil {
		return nil, errors.New("[DECODE] " + errDecoding.Error())
	}

	return message, nil
}

func encode(message interface{}, ch channel) ([]byte, error) {

	// Get schema from ch
	schema, errGetSchema := getSchema(ch)
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	// Get codec from schema
	codec, errCreateCodec := getCodec(*schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	messageEncoded, errParseAvro := codec.BinaryFromNative(nil, message)
	if errParseAvro != nil {
		return nil, errors.New("[ENCODE] " + errParseAvro.Error())
	}

	// // We append 5 extra bytes to the beginning of the messages
	// // in order to confluent avro serializer be able to deserrialize
	// // messages in other languages as python.
	// // see: https://docs.confluent.io/current/schema-registry/serializer-formatter.html#wire-format
	// // in section Wire Format.
	// magicByte := []byte{0}
	// schemaBytes := make([]byte, 4)
	// binary.BigEndian.PutUint32(schemaBytes, uint32(schema.ID))
	// wireBytes := append(magicByte, schemaBytes...)
	// messageEncoded = append(wireBytes, messageEncoded...)

	return messageEncoded, nil
}
