# Simple dockerfile for compiling flatbuffers.
# Author: Oliver 'kfsone' Smith <oliver@kfs.org>
# All rights lined up against a wall and mooned.
FROM alpine

RUN mkdir /gom && apk update && apk add make flatbuffers

VOLUME  ['/gom']
WORKDIR /gom

ENTRYPOINT [ "make", "flatc" ]

