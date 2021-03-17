package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/terra-project/mantle-sdk/serdes"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

func byteArrayToByteArray() mapstructure.DecodeHookFuncKind {
	return func(
		f reflect.Kind,
		t reflect.Kind,
		data interface{},
	) (interface{}, error) {

		return data, nil
	}
}

// TODO: remove me
func UnmarshalInternalQueryResult(result *graphql.Result, target interface{}) error {
	targetValue := reflect.Indirect(reflect.ValueOf(target))
	data, isDataMap := result.Data.(map[string]interface{})
	if !isDataMap {
		return fmt.Errorf("internal gql resulted in non-map data")
	}

	for key, pb := range data {
		targetField := targetValue.FieldByName(key)
		targetCache := reflect.New(targetField.Type())

		switch d := pb.(type) {
		case []byte:
			// amino or serdes bytes
			if d[0] == 196 {
				if err := json.Unmarshal(d[2:], targetCache.Interface()); err != nil {
					return err
				}
			} else {
				if unpackErr := serdes.Deserialize(targetField.Type(), d, targetCache.Interface()); unpackErr != nil {
					return unpackErr
				}
			}
			targetField.Set(targetCache.Elem())

		default:
			// TODO: FIX ME
			// map -> struct
			// mapstructure.ComposeDecodeHookFunc(byteArrayToByteArray)

			decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					byteArrayToByteArray(),
					mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05 +0000 UTC"),
				),
				WeaklyTypedInput: false,
				Result:           targetCache.Interface(),
				TagName:          "",
			})

			if decoderErr != nil {
				return decoderErr
			}
			if err := decoder.Decode(d); err != nil {
				return err
			}
			targetField.Set(targetCache.Elem())
		}
	}

	return nil
}

type Thunk func() (interface{}, error)
type ThunkResult struct {
	data interface{}
	err  error
}

func CreateParallel(len int) *parallelExecutionContext {
	wg := &sync.WaitGroup{}
	wg.Add(len)
	return &parallelExecutionContext{
		RWMutex: sync.RWMutex{},
		idx:     0,
		wg:      wg,
		result:  make([]ParallelExecutionResult, len),
	}
}

type parallelExecutionContext struct {
	sync.RWMutex
	idx         int64
	wg          *sync.WaitGroup
	result      []ParallelExecutionResult
	errorExists bool
	done        bool
}

type ParallelExecutionFunc func() (interface{}, error)
type ParallelExecutionResult struct {
	Result interface{}
	Error  error
}

func (pec *parallelExecutionContext) Run(f ParallelExecutionFunc) {
	if pec.done {
		panic("cannot add more runners. parallel execution is already done.")
	}

	i := pec.idx
	pec.idx = pec.idx + 1

	// run goroutine
	go func() {
		defer pec.wg.Done()

		r, e := f()

		pec.RWMutex.Lock()
		var result ParallelExecutionResult
		if e != nil {
			result = ParallelExecutionResult{Error: e}
			pec.errorExists = true
		} else {
			result = ParallelExecutionResult{Result: r}
		}

		pec.result[i] = result
		pec.RWMutex.Unlock()
	}()
}

func (pec *parallelExecutionContext) HasErrors() bool {
	return pec.errorExists
}

func (pec *parallelExecutionContext) Sync() []ParallelExecutionResult {
	pec.done = true
	pec.wg.Wait()

	return pec.result
}

func CreateRemoteMantleRequest(mantleEndpoint string, graphqlQuery []byte) []byte {
	request := new(struct {
		Query string `json:"query"`
	})

	request.Query = string(graphqlQuery)
	requestJson, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	res, err := http.Post(mantleEndpoint, "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		time.Sleep(200 * time.Millisecond)
		return CreateRemoteMantleRequest(mantleEndpoint, graphqlQuery)
	}

	gqlResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return gqlResponse
}

func buildSchema(schemabuilders ...SchemaBuilder) graphql.Schema {
	rootFields := &graphql.Fields{}

	for _, builder := range schemabuilders {
		err := builder(rootFields)

		if err != nil {
			panic(err)
		}
	}

	log.Println("[graphql] schemas", *rootFields)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: *rootFields,
		}),
	})

	if err != nil {
		panic(err)
	}

	return schema
}
