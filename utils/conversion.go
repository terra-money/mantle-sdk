package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/graphql-go/graphql"
	"math"
	"reflect"
)

var goTypeToGraphqlType = map[string]graphql.Type{
	"string":     graphql.String,
	"rune":       graphql.Int,
	"int":        graphql.Int,
	"byte":       graphql.Int,
	"int8":       graphql.Int,
	"int16":      graphql.Int,
	"int32":      graphql.Int,
	"int64":      graphql.Int,
	"uint8":      graphql.Int,
	"uint16":     graphql.Int,
	"uint32":     graphql.Int,
	"uint64":     graphql.Int,
	"bool":       graphql.Boolean,
	"float32":    graphql.Float,
	"float64":    graphql.Float,
	"complex64":  graphql.Float,
	"complex128": graphql.Float,
}

func GetGraphQLType(k reflect.Kind) graphql.Type {
	return goTypeToGraphqlType[k.String()]
}

func GetType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return GetType(t.Elem())
	default:
		return t
	}
}

func GetValue(v reflect.Value) reflect.Value {
	if v.Type().Kind() == reflect.Ptr {
		return GetValue(v.Elem())
	}
	return v
}

func IntToUintLexicographic(data int64) uint64 {
	if data >= 0 {
		return uint64(data) + uint64(math.MaxInt64)
	} else {
		return uint64(data + math.MaxInt64)
	}
}

func LeToBe(v uint64) []byte {
	be := make([]byte, 8)
	binary.BigEndian.PutUint64(be, v)
	return be
}

func BeToLe(v uint64) []byte {
	le := make([]byte, 8)
	binary.LittleEndian.PutUint64(le, v)
	return le
}

func ConvertToLexicographicBytes(data interface{}) ([]byte, error) {
	switch data.(type) {
	case string:
		return []byte(data.(string)), nil
	case uint:
		return LeToBe(uint64(data.(uint))), nil
	case uint8:
		return LeToBe(uint64(data.(uint8))), nil
	case uint16:
		return LeToBe(uint64(data.(uint16))), nil
	case uint32:
		return LeToBe(uint64(data.(uint32))), nil
	case uint64:
		return LeToBe(data.(uint64)), nil
	case int:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	case int8:
		return LeToBe(IntToUintLexicographic(int64(data.(int8)))), nil
	case int16:
		return LeToBe(IntToUintLexicographic(int64(data.(int16)))), nil
	case int32:
		return LeToBe(IntToUintLexicographic(int64(data.(int32)))), nil
	case int64:
		return LeToBe(IntToUintLexicographic(data.(int64))), nil
	default:
		return nil, fmt.Errorf("this type of data is disallowed for indexing: %s", reflect.TypeOf(data).Name())
	}
}

func ConvertToIndexValueToCorrectType(indexType reflect.Type, data interface{}) ([]byte, error) {
	switch indexType.Kind() {
	case reflect.String:
		return []byte(data.(string)), nil
	case reflect.Uint:
		return LeToBe(uint64(data.(int))), nil
	case reflect.Uint8:
		return LeToBe(uint64(data.(int))), nil
	case reflect.Uint16:
		return LeToBe(uint64(data.(int))), nil
	case reflect.Uint32:
		return LeToBe(uint64(data.(int))), nil
	case reflect.Uint64:
		return LeToBe(uint64(data.(int))), nil
	case reflect.Int:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	case reflect.Int8:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	case reflect.Int16:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	case reflect.Int32:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	case reflect.Int64:
		return LeToBe(IntToUintLexicographic(int64(data.(int)))), nil
	default:
		return nil, fmt.Errorf("this type of data is disallowed for indexing: %s", reflect.TypeOf(data).Name())
	}
}
