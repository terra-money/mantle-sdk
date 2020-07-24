package graph

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
)

// TODO: Make a better version of this or scrap
func UnmarshalGraphQLResult(result *graphql.Result, target interface{}) error {
	res, err := json.Marshal(result.Data)

	if err != nil {
		return err
	}

	err = json.Unmarshal(res, target)
	if err != nil {
		return err
	}

	return nil
}
