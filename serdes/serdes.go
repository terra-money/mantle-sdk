package serdes

import (
	"github.com/terra-project/core/app"
	"github.com/terra-money/mantle-sdk/types"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

type serializer func(data interface{}) ([]byte, error)
type deserializer func(data []byte, target interface{}) error

type serdes struct {
	serialize   serializer
	deserialize deserializer
}

var amino = app.MakeCodec()
var customSerializedModels = map[reflect.Type]serdes{
	reflect.TypeOf((*types.StdTx)(nil)).Elem(): {
		serialize: func(data interface{}) ([]byte, error) {
			return amino.MarshalBinaryLengthPrefixed(data)
		},
		deserialize: func(data []byte, target interface{}) error {
			return amino.UnmarshalBinaryLengthPrefixed(data, target)
		},
	},
}

func Serialize(t reflect.Type, data interface{}) ([]byte, error) {
	if customSerializer, ok := customSerializedModels[t]; ok {
		return customSerializer.serialize(data)
	}

	return DefaultMsgpackSerializer(data)
}

func Deserialize(t reflect.Type, data []byte, target interface{}) error {
	if customSerializer, ok := customSerializedModels[t]; ok {
		return customSerializer.deserialize(data, target)
	}

	return DefaultMsgpackDeserializer(data, target)
}

func DefaultMsgpackSerializer(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

func DefaultMsgpackDeserializer(data []byte, target interface{}) error {
	return msgpack.Unmarshal(data, target)
}
