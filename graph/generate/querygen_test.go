package generate

import (
	"fmt"
	"testing"

	types "github.com/terra-project/mantle/types"
)

type TestQueryStruct struct {
	TopLevelQuery string
	PlainQuery    struct {
		NestedField1         string
		NestedField2         int
		NestedCompositeField struct {
			HelloWorld string
		}
	}
	PlainQueryWithArguments struct {
		NestedField1         string
		NestedField2         int
		NestedCompositeField struct {
			HelloWorld string
		}
	} `mantle:"(filter: $isFieldPresent)"`
	AliasedQuery struct {
		NestedField1 string
	} `mantle:"PlainQuery"`
	AliasedQueryWithArguments struct {
		NestedField1 string `mantle:"(subfilter: $subfilter)"`
		NestedField2 int
	} `mantle:"PlainQuery(filter: $filter)"`
	BlockNow types.Block
	Blocks   []types.Block
}

func TestGenerateQuery(t *testing.T) {
	result := TestQueryStruct{}
	q := GenerateQuery(&result)
	fmt.Println(q)
}
