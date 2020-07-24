package utils

import (
	"io/ioutil"

	types "github.com/terra-project/mantle/types"
	yaml "gopkg.in/yaml.v2"
)

func GetTerraSwaggerDefinitions() types.Swagger {
	var swagger = types.Swagger{}

	swaggerFile, err := ioutil.ReadFile("./utils/swagger.yaml")

	if err != nil {
		panic("Loading swagger file failed")
	}

	if err := yaml.Unmarshal(swaggerFile, &swagger); err != nil {
		panic(err)
	}

	return swagger
}
