package schemabuilders

import (
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/generate"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
)

func CreateModelSchemaBuilder(kvindexMap kvindex.KVIndexMap, models ...types.Model) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		// handle module registration
		for _, model := range models {
			model = utils.GetType(model)
			modelName := model.Name()

			fieldConfig, err := generate.GenerateGraphResolver(model)
			if err != nil {
				return err
			}

			if fieldConfig == nil {
				continue
			}

			kvIndex, kvIndexExists := kvindexMap[modelName]
			if kvIndexExists {
				fieldConfig.Args = generate.GenerateArgument(kvIndex)
			}

			(*fields)[modelName] = fieldConfig

			// list (single entity is overwritten in case of slice model)
			listFieldConfig, err := generate.GenerateListGraphResolver(model, fieldConfig)
			if err != nil {
				return err
			}

			if listFieldConfig == nil {
				continue
			}

			(*fields)[listFieldConfig.Name] = listFieldConfig
		}

		return nil
	}
}
