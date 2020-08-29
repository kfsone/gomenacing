# Simple dockerfile for compiling flatbuffers.
# Author: Oliver 'kfsone' Smith <oliver@kfs.org>
# All rights lined up against a wall and mooned.
FROM alpine

RUN mkdir /ettudata && apk update && apk add make flatbuffers

VOLUME  ['/ettudata']
WORKDIR /ettudata

ENTRYPOINT [ "make", "flatc" ]

