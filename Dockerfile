# build-env stage

# librdkafka base image

FROM golang:alpine AS kafka_build

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git gcc g++ make bash pkgconfig
    # apk add --no-cache git bash libc-dev alpine-sdk

RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install
