package schemabuilders

import (
	"reflect"
	"sync"

	"github.com/go-openapi/strfmt"
	"github.com/graphql-go/graphql"
	terra "github.com/terra-project/core/app"
	compatlocalclient "github.com/terra-project/mantle-compatibility/localclient"
	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/graph/schemabuilders/abcistub"
	lcd "github.com/terra-project/mantle-sdk/lcd/client"
)

func CreateABCIStubSchemaBuilder(app *terra.TerraApp, queryMtx *sync.Mutex) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		localClient := compatlocalclient.NewLocalClient(app, queryMtx)
		stubTransport, err := abcistub.NewABCIStubTransport(localClient)
		if err != nil {
			return err
		}

		cli := lcd.New(stubTransport, strfmt.Default)
		cliv := reflect.ValueOf(cli).Elem()

		for i := 0; i < cliv.NumField(); i++ {
			vf := cliv.Field(i)
			vt := vf.Type()

			for j := 0; j < vf.NumMethod(); j++ {
				clientFunc := vf.Method(j)
				clientFuncName := vt.Method(j).Name
				clientFuncType := clientFunc.Type()

				if clientFuncName[:3] != "Get" {
					continue
				}

				fieldConfig, err := abcistub.RegisterABCIQueriers(clientFunc, clientFuncName, clientFuncType)
				if err != nil {
					return err
				}

				if fieldConfig == nil {
					continue
				}

				canonicalName := clientFuncName[3:]
				(*fields)[canonicalName] = fieldConfig
			}
		}

		return nil
	}
}
