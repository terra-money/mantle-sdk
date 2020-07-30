package abcistub

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/strfmt"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/lcd/lcd"
	"github.com/terra-project/mantle/utils"
)

type contextKey int

const varsKey contextKey = iota

func RegisterABCIQueriers(fields *graphql.Fields, localClient LocalClient) error {

	// initialize lcd
	stubTransport, err := NewABCIStubTransport(localClient)
	if err != nil {
		return err
	}

	cli := lcd.New(stubTransport, strfmt.Default)

	v := reflect.ValueOf(cli).Elem()

	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		vt := vf.Type()

		for j := 0; j < vf.NumMethod(); j++ {
			f2 := vf.Method(j)

			fnName := vt.Method(j).Name
			fnType := f2.Type()

			// do NOT handle non-get functions
			if fnName[:3] != "Get" {
				continue
			}

			canonicalName := fnName[3:]
			argumentsType := buildArguments(fnType)

			out, exists := fnType.Out(0).Elem().FieldByName("Payload")
			if !exists {
				continue
			}

			responseType := buildResponseType(out.Type, out.Name, fnName)

			gqlField := graphql.Field{
				Name: canonicalName,
				Type: responseType,
				Args: argumentsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := utils.GetType(fnType.In(0))

					// map[string]interface{} -> struct (as defined in the handler function)
					argsStruct := reflect.New(args)
					for key, value := range p.Args {
						argsStruct.Elem().FieldByName(key).Set(reflect.ValueOf(value))
					}

					result := f2.Call([]reflect.Value{argsStruct})
					ret, err := result[0], result[1]

					if !err.IsNil() {
						return nil, err.Interface().(error)
					}

					return utils.GetValue(ret.Elem().FieldByName("Payload")).Interface(), nil
				},
			}

			(*fields)[canonicalName] = &gqlField
		}
	}

	return nil
}

func buildResponseType(t reflect.Type, tName string, parentName string) graphql.Output {
	kind := t.Kind()

	switch kind {
	// in case of struct,
	case reflect.Struct:
		fields := graphql.Fields{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			fields[field.Name] = &graphql.Field{
				Name: field.Name,
				Type: buildResponseType(field.Type, field.Name, tName),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return utils.GetValue(reflect.ValueOf(p.Source)).FieldByName(field.Name).Interface(), nil
				},
			}
		}

		return graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%s%s", parentName, tName),
			Fields: fields,
		})

		// in case of ptr, take Elem() of the type and go deeper
	case reflect.Ptr:
		t := t.Elem()
		return buildResponseType(t, t.Name(), parentName)

		// in case of slice,
	case reflect.Slice, reflect.Array:
		t := t.Elem()
		return graphql.NewList(buildResponseType(t, t.Name(), parentName))

	default:
		return utils.GetGraphQLType(kind)
	}
}

func buildArguments(model reflect.Type) graphql.FieldConfigArgument {
	var filter = []string{"timeout", "Context", "HTTPClient"}
	var args = graphql.FieldConfigArgument{}

	argumentType := model.In(0).Elem()

	for i := 0; i < argumentType.NumField(); i++ {
		var skip = false
		field := argumentType.Field(i)

		// skip fields defined in filter
		for _, skippedField := range filter {
			if field.Name == skippedField {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		// if field is a pointer type, then it is an optional type
		var argumentScalar graphql.Input

		if field.Type.Kind() == reflect.Ptr {
			argumentScalar = utils.GetGraphQLType(field.Type.Elem().Kind())
		} else {
			argumentScalar = utils.GetGraphQLType(field.Type.Kind())
		}

		if argumentScalar == nil {
			panic("Unknown input type detected")
		}

		args[field.Name] = &graphql.ArgumentConfig{
			Type: argumentScalar,
		}

	}

	return args
}
