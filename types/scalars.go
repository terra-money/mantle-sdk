package types

import (
	"encoding/json"
)

type JSONScalar string
type JSONScalarCustomMarshaler func(originalData interface{}, data []byte) []byte

func NewJSONScalar(data interface{}, customMarshaller JSONScalarCustomMarshaler) JSONScalar {
	bz, bzErr := json.Marshal(data)
	if bzErr != nil {
		panic(bzErr)
	}

	if customMarshaller != nil {
		bz = customMarshaller(data, bz)
	}

	return JSONScalar(bz)
}
func (scalar *JSONScalar) IsJsonScalar() {}
