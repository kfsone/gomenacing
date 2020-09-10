package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type DataLoader struct {
	unmarshaler func([]byte) error
	loader      func() error
}

func (l *DataLoader) Load(data []byte) (err error) {
	if err = l.unmarshaler(data); err == nil {
		err = l.loader()
	}
	return
}

func NewDataLoader(unmarshaler func([]byte) error, callback func() error) *DataLoader {
	return &DataLoader{unmarshaler: unmarshaler, loader: callback}
}

func NewProtoLoader(into proto.Message, callback func() error) *DataLoader {
	return NewDataLoader(func(data []byte) error { return proto.Unmarshal(data, into) }, callback)
}

func NewJSONLoader(into interface{}, callback func() error) *DataLoader {
	return NewDataLoader(func(data []byte) error { return json.Unmarshal(data, into) }, callback)
}

func NewTypedDataLoader(name string, storage interface{}, callback func() error) (*DataLoader, error) {
	if storage == nil {
		return nil, errors.New("nil storage for loader")
	}
	switch name {
	case "proto":
		return NewProtoLoader(storage.(proto.Message), callback), nil
	case "gom":
		return NewProtoLoader(storage.(proto.Message), callback), nil

	case "json":
		return NewJSONLoader(storage, callback), nil

	default:
		return nil, fmt.Errorf("%w: unrecognized loader type: %s", ErrUnknownEntity, name)
	}
}
