package types

import (
	"encoding/json"
)

type JSONScalar string

func NewJSONScalar(data interface{}) JSONScalar {
	bz, bzErr := json.Marshal(data)
	if bzErr != nil {
		panic(bzErr)
	}
	return JSONScalar(bz)
}
func (scalar *JSONScalar) IsJsonScalar() {}
