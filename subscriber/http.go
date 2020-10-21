package subscriber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/terra-project/mantle-sdk/types"
)

type BlockGetter func(height interface{}) (*types.Block, error)

func GetBlock(endpoint string) (*types.Block, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	resbytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(struct {
		Result struct {
			Block json.RawMessage
		} `json:"result"`
	})

	if unmarshalErr := json.Unmarshal(resbytes, data); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	block := types.Block{}

	if err := json.Unmarshal(data.Result.Block, &block); err != nil {
		panic(err)
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
