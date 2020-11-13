package schemabuilders

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/db/kvindex"
	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/serdes"
	"github.com/terra-project/mantle-sdk/types"
	"github.com/terra-project/mantle-sdk/utils"
	"reflect"
)

// TODO: make me into a separate module
func CreateMantleStateSchemaBuilder(_ kvindex.KVIndexMap, _ ...types.Model) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		(*fields)["LastSyncedHeight"] = &graphql.Field{
			Name:              "LastSyncedHeight",
			Type:              graphql.Int,
			Args:              nil,
			Resolve:           func(p graphql.ResolveParams) (interface{}, error) {
				querier, ok := p.Context.Value(utils.QuerierKey).(querier.Querier)
				if !ok {
					panic(fmt.Errorf("querier is either cleared or not set, in Resolver Func for LastSyncedHeight"))
				}

				payload, err := querier.Get([]byte("LastSyncedHeight"))

				var height uint64
				serdesErr := serdes.Deserialize(reflect.TypeOf(uint64(0)), payload, &height)
				if serdesErr != nil {
					return nil, err
				}

				if err != nil {
					return nil, fmt.Errorf("LastSyncedHeight doesnt exist")
				}

				return height, nil
			},
			DeprecationReason: "",
			Description:       "Last known height in mantle. Height is updated only after committing all entities, therefore it is guaranteed to be in sync with other indexer results.",
		}


		return nil
	}
}
