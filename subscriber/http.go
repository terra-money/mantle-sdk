package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/terra-project/mantle/types"
	"io/ioutil"
	"net/http"
)

type BlockGetter func(height interface{}) (*types.Block, error)

func GetBlockLCD(endpoint string) (*types.Block, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	resbytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	temp := map[string]json.RawMessage{}
	if err := json.Unmarshal(resbytes, &temp); err != nil {
		return nil, err
	}

	block := types.Block{}
	if err := json.Unmarshal(temp["block"], &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func CreateBlockGetterOffline() BlockGetter {
	return func(height interface{}) (*types.Block, error) {
		return nil, nil
	}
}

func forceString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	} else {
		return fmt.Sprintf("%s", val)
	}
}
