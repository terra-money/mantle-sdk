package abcistub

import (
	"fmt"
	"github.com/terra-project/mantle-sdk/graph/scalars"
	"reflect"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/utils"
)

type contextKey int

const varsKey contextKey = iota

func RegisterABCIQueriers(clientFunc reflect.Value, clientFuncName string, clientFuncType reflect.Type) (*graphql.Field, error) {
	canonicalName := clientFuncName[3:]
	argumentsType := buildArguments(clientFuncType)

	out, exists := clientFuncType.Out(0).Elem().FieldByName("Payload")
	if !exists {
		return nil, nil
	}

	responseType := buildResponseType(out.Type, out.Name, clientFuncName)

	return &graphql.Field{
		Name: canonicalName,
		Type: responseType,
		Args: argumentsType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			args := utils.GetType(clientFuncType.In(0))
			argsStruct := reflect.New(args)
			for key, value := range p.Args {
				argsStruct.Elem().FieldByName(key).Set(reflect.ValueOf(value))
			}

			result := clientFunc.Call([]reflect.Value{argsStruct})
			ret, err := result[0], result[1]

			if !err.IsNil() {
				return nil, err.Interface().(error)
			}

			return ret.Elem().FieldByName("Payload").Interface(), nil
		},
	}, nil
}

func buildResponseType(t reflect.Type, tName string, parentName string) graphql.Output {
	kind := t.Kind()

	switch kind {
	// in case of struct,
	case reflect.Struct:
		fields := graphql.Fields{}
		structName := fmt.Sprintf("%s%s", parentName, tName)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			var fieldType graphql.Output

			// see if this field should be implemented as scalar type
			scalar, isScalar := scalars.IsCosmosScalar(field.Type)
			if isScalar {
				fieldType = scalar
			} else {
				fieldType = buildResponseType(field.Type, field.Name, structName)

				// skip nil fields
				if fieldType == nil {
					continue
				}
			}

			fields[field.Name] = &graphql.Field{
				Name: field.Name,
				Type: fieldType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, isSourceValue := p.Source.(reflect.Value)
					if !isSourceValue {
						source = reflect.ValueOf(p.Source)
					}
					return reflect.Indirect(source).FieldByName(field.Name).Interface(), nil
				},
			}
		}

		return graphql.NewObject(graphql.ObjectConfig{
			Name:   structName,
			Fields: fields,
		})

	case reflect.Interface:
		return graphql.NewScalar(graphql.ScalarConfig{
			Name:         fmt.Sprintf("%s%s", parentName, tName),
			Serialize:    func(value interface{}) interface{} { return value },
			ParseValue:   func(value interface{}) interface{} { return value },
			ParseLiteral: func(valueAST ast.Value) interface{} { return nil },
		})

	// in case of ptr, take Elem() of the type and go deeper
	case reflect.Ptr:
		t := t.Elem()
		return buildResponseType(t, tName, parentName)

		// in case of slice,
	case reflect.Slice, reflect.Array:
		t := t.Elem()
		responseType := buildResponseType(t, tName, parentName)
		if responseType == nil {
			return nil
		}
		return graphql.NewList(responseType)

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
