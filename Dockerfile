# Simple dockerfile for compiling flatbuffers.
# Author: Oliver 'kfsone' Smith <oliver@kfs.org>
# All rights lined up against a wall and mooned.
FROM golang:1.14-alpine

RUN apk update \
		&& apk add git make protoc python3 && \
		go get -u google.golang.org/protobuf/cmd/protoc-gen-go && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go

VOLUME  ['/gom']
WORKDIR /gom

ENTRYPOINT [ "make", "protoc" ]

