package schemabuilders

import (
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/generate"
	"github.com/terra-project/mantle/types"
)

func CreateModelSchemaBuilder(models ...types.ModelType) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		// handle module registration
		for _, model := range models {
			model = utils.GetType(model)
			entityName := model.Name()
			fieldConfig, err := generate.GenerateGraphResolver(model)
			if err != nil {
				return err
			}

			if fieldConfig == nil {
				continue
			}

			fieldConfig.Args = generate.GenerateArgument(kvindex.NewKVIndex(model))
			(*fields)[entityName] = fieldConfig

			// list
			entityNamePlural := utils.Pluralize(entityName)
			listFieldConfig, err := generate.GenerateListGraphResolver(model, (*fields)[entityName])
			if err != nil {
				return err
			}

			if listFieldConfig == nil {
				continue
			}

			(*fields)[entityNamePlural] = listFieldConfig
		}

		return nil
	}
}
