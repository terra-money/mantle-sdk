package utils

import (
	"encoding/json"
	"fmt"
)

func DebugJSONPrint(data []byte, v interface{}) {
	json.Unmarshal(data, &v)
	byte, _ := json.MarshalIndent(v, "", "")
	fmt.Println(byte)
}
