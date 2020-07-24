package schemabuilders

import (
	"github.com/graphql-go/graphql"
	terra "github.com/terra-project/core/app"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/schemabuilders/abcistub"
)

func CreateABCIStubSchemaBuilder(app *terra.TerraApp) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		localClient := abcistub.NewLocalClient(app)
		return abcistub.RegisterABCIQueriers(fields, localClient)
	}
}
