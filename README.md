# Chimera client

# Overview

The clients provide a friendly interface that deal with read and write messages (en/coded in `Avro` schemas) on `Chimera` channels.

## What these clients can do?

- Read messages from channels;
- Write messages on channels;
- Get the `Avro` schema of a channel;
- Get the object with all environment variables;
- Log messages sent and received in the log channel.

# Environment Variables

The Clients API consists of two simple classes. A `Reader` and a `Writer` .  Both of these require all the environment variables set. These variables defines the necessary configuration that the Clients need to work properly.

- CHIMERA_NODE_ID
- CHIMERA_REGISTRY_URL
- KAFKA_BOOTSTRAP_SERVERS
- CHIMERA_NAMESPACE
- CHIMERA_INPUT_CHANNELS
- CHIMERA_OUTPUT_CHANNELS

## Allowed Channels

Let's talks about the `CHIMERA_INPUT_CHANNELS` and `CHIMERA_OUTPUT_CHANNELS` environment variables. When instantiating a pipeline, the user must define channels which the client is allowed to read and write messages, for example, consider the following `my_pipeline.yaml`:

```yaml

version: 1.0
namespace: my_namespace
name: pipeline_process_and_send
apps:
  - name: processing
    imagepath: image_path_for_processing:latest
    dockercontext: path_of_processing
    inputs:
        - pre_processed_data
        - raw_data_from_db
    outputs:
        - upload_data
  - name: uploader
    imagepath: image_path_for_uploader:latest
    context: path_of_uploader
    inputs:
        - upload_data
```

In this pipeline, there are two apps that will make use of the pipeline: **processing** and **uploader** app.

The **processing** service is only allowed to read messages from *pre_processed_data* and *raw_data_from_db* channels and allowed to write messages in the *upload_data* channel. 

The **uploader** service is allowed to read messages only from the *upload_data* channel and it is not allowed to write messages in any channel.

If the user tries to read/write messages from/to some channel different from these ones, the client will return an error saying that it's not allowed to perform the operation.

# Avro serialization

We opted to save the massages in channels using the `Avro` coding. Why? Well, `Kafka` platform save message as byte arrays, which is basically strings, so we choose to save the messages using the `JSON` format, because objects are well defined and converted to strings. 

But, `JSON` format requires pairs of key:value and we don't need to pass the keys everytime if we know that it always be the same. So we define `Avro` schemas to remember this keys and en/conding the messages, saving a lot of storage passing only the values. 

See more: 

[https://avro.apache.org/docs/current/index.html](https://avro.apache.org/docs/current/index.html)

# Usage
https://godoc.org/github.com/insprbr/client-golang

Reader and writer examples

## Reader

For now, `Reader` fetch a single message, listening in the input channels or to a specific channel, when a message comes in.  We plan for future versions implement the option to read messages in batches.

Just after the client read message, it is decoded according to the specified schema for the channel and returns an object with the original message.


```go
// New Reader instance
reader, errNewReader := client.NewReader()
if errNewReader != nil {
    log.Fatal("[INSTANTIATING NEW READER] " + errNewReader.Error())
}

// Reading message
msgChannel, msg, errReadMessage := reader.ReadMessage()
log.Println(*msgChannel)
if errReadMessage != nil {
    log.Fatal("[READ MESSAGE] " + errReadMessage.Error())
}
```

## Writer

The Writer consist in a single functionality to write a message in a desired, allowed `Chimera` channel.

```go
// New Writer instance
writer, err = client.NewWriter()
if err != nil {
    log.Panicln(err)
}

// message
user_info_interface = map[string]interface{}{
    "name":      "user_name",
    "age":       18,
    "Gender"     "male"
}

// Writing message on channel target_channel
err := writer.WriteMessage(user_info_interface, "target_channel")
if err != nil {
    log.Panicln(err)
}
```
