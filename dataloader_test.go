package main

import (
	"encoding/json"
	"fmt"
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestDataLoader_Load(t *testing.T) {
	t.Run("error from Marshaler", func(t *testing.T) {
		loads := 0
		loader := DataLoader{
			unmarshaler: func(value []byte) error {
				return fmt.Errorf("unmarshal: |%s|", string(value))
			},
			loader: func() error {
				loads++
				return nil
			},
		}
		err := loader.Load([]byte("VALUE"))
		if assert.Error(t, err) {
			assert.Equal(t, 0, loads)
			assert.Equal(t, "unmarshal: |VALUE|", err.Error())
		}
	})

	t.Run("error from Loader", func(t *testing.T) {
		loads := 0
		var data []byte
		loader := DataLoader{
			unmarshaler: func(value []byte) error {
				data = value
				return nil
			},
			loader: func() error {
				loads++
				return fmt.Errorf("load: |%s|", string(data))
			},
		}
		err := loader.Load([]byte("VALUE"))
		if assert.Error(t, err) {
			assert.Equal(t, 1, loads)
			assert.Equal(t, "load: |VALUE|", err.Error())
		}
	})

	t.Run("pass through", func(t *testing.T) {
		var data []byte
		var received string
		loader := DataLoader{
			unmarshaler: func(value []byte) error {
				data = value
				return nil
			},
			loader: func() error {
				received += string(data) + "\n"
				return nil
			},
		}
		if assert.Nil(t, loader.Load([]byte("hello"))) {
			if assert.Nil(t, loader.Load([]byte("world"))) {
				assert.Equal(t, "hello\nworld\n", received)
			}
		}
	})
}

func TestNewDataLoader(t *testing.T) {
	unmarshals := 0
	unmarshaler := func([]byte) error { unmarshals++; return nil }
	loads := 0
	callback := func() error { loads++; return nil }
	loader := NewDataLoader(unmarshaler, callback)
	if assert.NotNil(t, loader) {
		assert.Nil(t, loader.unmarshaler([]byte("")))
		assert.Equal(t, 1, unmarshals)
		assert.Equal(t, 0, loads)

		assert.Nil(t, loader.loader())
		assert.Equal(t, 1, unmarshals)
		assert.Equal(t, 1, loads)
	}
}

func TestNewProtoLoader(t *testing.T) {
	into := &gomschema.Header{}
	loads := 0
	callback := func() error { loads++; return nil }
	loader := NewProtoLoader(into, callback)
	if assert.NotNil(t, loader) {
		assert.NotNil(t, loader.unmarshaler)
		assert.NotNil(t, loader.loader)
		message, err := proto.Marshal(&gomschema.Header{
			HeaderType: gomschema.Header_CFacility,
			Sizes:      []uint32{1, 10, 1},
			Source:     "testing",
			Userdata:   nil,
		})
		require.Nil(t, err)
		require.NotNil(t, message)

		if assert.Nil(t, loader.Load(message)) {
			assert.Equal(t, 1, loads)
			assert.Equal(t, gomschema.Header_CFacility, into.HeaderType)
			assert.Equal(t, []uint32{1, 10, 1}, into.Sizes)
			assert.Equal(t, "testing", into.Source)
			assert.Nil(t, into.Userdata)
		}
	}
}

func TestJSONLoader(t *testing.T) {
	type TestMessage struct {
		ID   int
		Name string
	}

	loads := 0
	callback := func() error { loads++; return nil }
	into := &TestMessage{}
	loader := NewJSONLoader(into, callback)
	if assert.NotNil(t, loader) {
		assert.NotNil(t, loader.unmarshaler)
		assert.NotNil(t, loader.loader)
		message, err := json.Marshal(&TestMessage{42, "lt&e"})
		require.Nil(t, err)
		require.NotNil(t, message)

		if assert.Nil(t, loader.Load(message)) {
			assert.Equal(t, 1, loads)
			assert.Equal(t, 42, into.ID)
			assert.Equal(t, "lt&e", into.Name)
		}
	}
}

func TestNewTypedDataLoader(t *testing.T) {
	loader, err := NewTypedDataLoader("nil", nil, func() error { return nil })
	if assert.Error(t, err) {
		assert.Nil(t, loader)
	}

	loader, err = NewTypedDataLoader("invalid", nil, func() error { return nil })
	if assert.Error(t, err) {
		assert.Nil(t, loader)
	}

	protoData, err := proto.Marshal(&gomschema.Header{})
	require.NotNil(t, protoData)
	require.Nil(t, err)

	jsonData, err := json.Marshal("hello")
	require.NotNil(t, jsonData)
	require.Nil(t, err)

	// Make sure the proto loader can decode protobuffs
	loader, err = NewTypedDataLoader("proto", &gomschema.Header{}, func() error { return nil })
	if assert.Nil(t, err) {
		assert.NotNil(t, loader)
		assert.Nil(t, loader.Load(protoData))
		assert.Error(t, loader.Load(jsonData))
	}

	loader, err = NewTypedDataLoader("proto", &gomschema.Header{}, func() error { return nil })
	if assert.Nil(t, err) {
		assert.NotNil(t, loader)
		assert.Nil(t, loader.Load(protoData))
		assert.Error(t, loader.Load(jsonData))
	}

	// Make sure the json loader can decode protos
	var test string
	loader, err = NewTypedDataLoader("json", &test, func() error { return nil })
	if assert.Nil(t, err) {
		assert.NotNil(t, loader)
		assert.Error(t, loader.Load(protoData))
		assert.Nil(t, loader.Load(jsonData))
	}
}
