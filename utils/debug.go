package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func DebugJSONPrint(data []byte, v interface{}) {
	json.Unmarshal(data, &v)
	byte, _ := json.MarshalIndent(v, "", "")
	fmt.Println(byte)
}

func DebugStoreKey(key []byte) {
	// split seq delimiter
	seq := strings.Split(string(key), string(DocumentSeqDelimiter))

	// split index value delimiter
	iv := strings.Split(seq[0], string(DocumentIndexValueDelimiter))

	//
	ik := strings.Split(iv[0], string(DocumentIndexKeyDelimiter))

	if len(ik) > 1 {
		fmt.Printf(
			"document(%s) indexName(%s) indexValue(%v) seq(%v)\n",
			ik[0],
			ik[1],
			[]byte(iv[1]),
			[]byte(seq[1]),
		)
	} else {
		fmt.Printf(
			"document(%s) seq(%v)\n",
			ik[0],
			[]byte(seq[1]),
		)
	}

}
