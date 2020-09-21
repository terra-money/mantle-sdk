package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MustUnmarshal panics if json unmarshal fails.
// Use with care.
func MustUnmarshal(data []byte, target interface{}) {
	fmt.Println(string(data))
	if err := json.Unmarshal(data, target); err != nil {
		panic(err)
	}
}

// HasKey is 1 level deep key finder.
// returns false if key is not found.
func IsJSONKeyPresent(data interface{}, key string) bool {
	return !reflect.Indirect(reflect.ValueOf(data)).FieldByName(key).IsZero()
}
