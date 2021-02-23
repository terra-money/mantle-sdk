package proxy_resolver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OfType struct {
	Kind   string  `json:"kind"`
	Name   string  `json:"name"`
	OfType *OfType `json:"ofType"`
}
type Type struct {
	Kind   string      `json:"kind"`
	Name   interface{} `json:"name"`
	OfType *OfType     `json:"ofType"`
}
type Args struct {
	DefaultValue interface{} `json:"defaultValue"`
	Description  string      `json:"description"`
	Name         string      `json:"name"`
	Type         Type        `json:"type"`
}
type Directives struct {
	Args        []Args   `json:"args"`
	Description string   `json:"description"`
	Locations   []string `json:"locations"`
	Name        string   `json:"name"`
}
type QueryType struct {
	Name string `json:"name"`
}

type Argument struct {
	DefaultValue interface{} `json:"defaultValue"`
	Description  string      `json:"description"`
	Name         string      `json:"name"`
	Type         *Input      `json:"type"`
}
type Input struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	OfType *Input `json:"ofType"`
}

type Field struct {
	Args              []Argument `json:"args"`
	DeprecationReason string     `json:"deprecationReason"`
	Description       string     `json:"description"`
	IsDeprecated      bool       `json:"isDeprecated"`
	Name              string     `json:"name"`
	Type              Type       `json:"type"`
}
type EnumValue struct {
	DeprecationReason string `json:"deprecationReason"`
	Description       string `json:"description"`
	IsDeprecated      bool   `json:"isDeprecated"`
	Name              string `json:"name"`
}
type TypeDescriptor struct {
	Description   string        `json:"description"`
	EnumValues    []EnumValue   `json:"enumValues"`
	Fields        []Field       `json:"fields"`
	InputFields   interface{}   `json:"inputFields"`
	Interfaces    []interface{} `json:"interfaces"`
	Kind          string        `json:"kind"`
	Name          string        `json:"name"`
	PossibleTypes interface{}   `json:"possibleTypes"`
}
type Schema struct {
	Directives       []Directives     `json:"directives"`
	MutationType     interface{}      `json:"mutationType"`
	QueryType        QueryType        `json:"queryType"`
	SubscriptionType interface{}      `json:"subscriptionType"`
	Types            []TypeDescriptor `json:"types"`
}

var introspectionQuery = "query IntrospectionQuery {" +
	"  __schema {" +
	"    queryType {" +
	"      name" +
	"    }" +
	"    mutationType {" +
	"      name" +
	"    }" +
	"    subscriptionType {" +
	"      name" +
	"    }" +
	"    types {" +
	"      ...FullType" +
	"    }" +
	"    directives {" +
	"      name" +
	"      description" +
	"      locations" +
	"      args {" +
	"        ...InputValue" +
	"      }" +
	"    }" +
	"  }" +
	"}" +
	"" +
	"fragment FullType on __Type {" +
	"  kind" +
	"  name" +
	"  description" +
	"  fields(includeDeprecated: true) {" +
	"    name" +
	"    description" +
	"    args {" +
	"      ...InputValue" +
	"    }" +
	"    type {" +
	"      ...TypeRef" +
	"    }" +
	"    isDeprecated" +
	"    deprecationReason" +
	"  }" +
	"  inputFields {" +
	"    ...InputValue" +
	"  }" +
	"  interfaces {" +
	"    ...TypeRef" +
	"  }" +
	"  enumValues(includeDeprecated: true) {" +
	"    name" +
	"    description" +
	"    isDeprecated" +
	"    deprecationReason" +
	"  }" +
	"  possibleTypes {" +
	"    ...TypeRef" +
	"  }" +
	"}" +
	"" +
	"fragment InputValue on __InputValue {" +
	"  name" +
	"  description" +
	"  type {" +
	"    ...TypeRef" +
	"  }" +
	"  defaultValue" +
	"}" +
	"" +
	"fragment TypeRef on __Type {" +
	"  kind" +
	"  name" +
	"  ofType {" +
	"    kind" +
	"    name" +
	"    ofType {" +
	"      kind" +
	"      name" +
	"      ofType {" +
	"        kind" +
	"        name" +
	"        ofType {" +
	"          kind" +
	"          name" +
	"          ofType {" +
	"            kind" +
	"            name" +
	"            ofType {" +
	"              kind" +
	"              name" +
	"              ofType {" +
	"                kind" +
	"                name" +
	"              }" +
	"            }" +
	"          }" +
	"        }" +
	"      }" +
	"    }" +
	"  }" +
	"}"

func NewIntrospection(url string) Schema {
	gqlbody := new(struct {
		OperationName string    `json:"operationName"`
		Query         string    `json:"query"`
		Variables     *struct{} `json:"variables"`
	})

	gqlbody.OperationName = "IntrospectionQuery"
	gqlbody.Query = introspectionQuery
	gqlbody.Variables = nil
	gqlbodyBz, _ := json.Marshal(gqlbody)

	response, httpErr := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(gqlbodyBz),
	)

	if httpErr != nil {
		invariant(httpErr)
	}

	bz, bzErr := ioutil.ReadAll(response.Body)
	if bzErr != nil {
		invariant(bzErr)
	}

	introspection := new(struct {
		Data struct {
			Schema Schema `json:"__schema"`
		} `json:"data"`
	})

	if err := json.Unmarshal(bz, &introspection); err != nil {
		invariant(err)
	}

	return introspection.Data.Schema
}

func invariant(e error) {
	panic(fmt.Errorf("remote mantle introspection failed, %s", e.Error()))
}
