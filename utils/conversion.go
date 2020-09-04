package utils

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"

	"github.com/graphql-go/graphql"
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
	if t.Kind() == reflect.Ptr {
		return GetType(t.Elem())
	} else {
		return t
	}
}

func GetValue(v reflect.Value) reflect.Value {
	if v.Type().Kind() == reflect.Ptr {
		return GetValue(v.Elem())
	}
	return v
}

func GetUint64FromWhatever(v interface{}) (uint64, error) {
	k := reflect.TypeOf(v).Kind()

	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret, _ := strconv.ParseUint(fmt.Sprintf("%d", v), 10, 64)
		return ret, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret, _ := strconv.ParseUint(fmt.Sprintf("%d", v), 10, 64)
		return ret, nil
	case reflect.String:
		ret, _ := strconv.ParseUint(v.(string), 10, 64)
		return ret, nil
	}

	return 0, fmt.Errorf("Value is not convertible to uint64")
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
