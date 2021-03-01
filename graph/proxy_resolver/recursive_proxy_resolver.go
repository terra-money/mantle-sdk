package proxy_resolver

import (
	"bytes"
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	"strconv"
	"sync"
)

type ProxyResolverContext struct {
	name         string
	alias        string
	isRoot       bool
	mtx          *sync.Mutex
	arguments    map[string]ast.Value
	variables    map[string]interface{}
	value        map[string]interface{}
	error        error
	responseCb   ProxyResolverResponseCallback
	parent       *ProxyResolverContext
	subtree      []*ProxyResolverContext
	subtreeNames []string
}

type ProxyResolverResponseCallback func(query []byte) (map[string]interface{}, error)

var ProxyResolverContextKey = "proxy-resolver-context"

func NewProxyResolverContext(responseCb ProxyResolverResponseCallback) *ProxyResolverContext {
	return &ProxyResolverContext{
		name:      "root",
		mtx:       new(sync.Mutex),
		isRoot:    true,
		arguments: nil,
		variables: nil, // only set in root
		value:     nil,
		error:     nil,
		parent:    nil,
		// root
		responseCb:   responseCb,
		subtree:      make([]*ProxyResolverContext, 0),
		subtreeNames: make([]string, 0),
	}
}

// TODO: make subtree a separate struct, implementing the same interface
func (prc *ProxyResolverContext) CreateSubtree(typeName string, arguments map[string]ast.Value) *ProxyResolverContext {
	sprc := &ProxyResolverContext{
		name:         typeName,
		isRoot:       false,
		arguments:    arguments,
		variables:    nil, // never set for subtree
		mtx:          nil, // never set for subtree
		value:        nil,
		error:        nil,
		responseCb:   nil,
		parent:       prc,
		subtree:      make([]*ProxyResolverContext, 0),
		subtreeNames: make([]string, 0),
	}

	prc.subtree = append(prc.subtree, sprc)
	prc.subtreeNames = append(prc.subtreeNames, typeName)

	return sprc
}

func (prc *ProxyResolverContext) WithAlias(alias string) *ProxyResolverContext {
	prc.alias = alias
	return prc
}

func (prc *ProxyResolverContext) WithGraphQLVariables(variables map[string]interface{}) *ProxyResolverContext {
	if !prc.isRoot {
		panic("only root prc can contain arguments map")
	}

	prc.variables = variables
	return prc
}

func (prc *ProxyResolverContext) ResolveRoot() (map[string]interface{}, error) {
	if !prc.isRoot {
		return nil, errNoRoot
	}

	// set mutex here -- whichever root proxy context that enters here would put a lock,
	// basically barring others from executing N queries
	prc.mtx.Lock()
	defer prc.mtx.Unlock()

	// subsequent calls to ResolveRoot would serve from cache
	if prc.value != nil || len(prc.value) != 0 {
		return prc.value, prc.error
	}

	if prc.error != nil {
		return nil, prc.error
	}

	// reconstruct query using prc
	query := reconstructRootQuery(prc)

	if value, err := prc.responseCb(query.Bytes()); err == nil {
		prc.value = value
	} else {
		prc.error = err
	}

	return prc.value, prc.error
}

func reconstructRootQuery(prc *ProxyResolverContext) *bytes.Buffer {
	// query
	queryBuf := new(bytes.Buffer)

	// keep a copy of root variables -- only root prc has it
	variables := prc.variables

	// turn everything into query
	reconstructSubselectionQuery(queryBuf, prc, variables)

	return queryBuf
}

func reconstructSubselectionQuery(queryBuf *bytes.Buffer, prc *ProxyResolverContext, variables map[string]interface{}) {
	// write arguments part
	reconstructSubselectionArgument(queryBuf, prc.arguments, variables)

	// write subselection parts
	hasSubselection := len(prc.subtree) != 0

	if hasSubselection {
		queryBuf.Write([]byte("{"))
	}
	for subtreeIdx, subtreeItem := range prc.subtree {
		subtreeName := prc.subtreeNames[subtreeIdx]
		subtreePrc := prc.subtree[subtreeIdx]

		// write alias
		if subtreePrc.alias != "" {
			queryBuf.Write([]byte(subtreePrc.alias))
			queryBuf.Write([]byte(":"))
		}

		queryBuf.Write([]byte(subtreeName))
		reconstructSubselectionQuery(queryBuf, subtreeItem, variables)
		queryBuf.Write([]byte(","))
	}
	if hasSubselection {
		queryBuf.Write([]byte("}"))
	}
}

func reconstructSubselectionArgument(qb *bytes.Buffer, arguments map[string]ast.Value, variables map[string]interface{}) {
	if arguments == nil || len(arguments) == 0 {
		return
	}

	qb.Write([]byte("("))
	for ak, av := range arguments {
		qb.Write([]byte(ak))
		qb.Write([]byte(":"))
		qb.Write(toCorrectGraphQLArgument(av, variables))
		qb.Write([]byte(","))
	}
	qb.Write([]byte(")"))
}

func toCorrectGraphQLArgument(argValue ast.Value, variables map[string]interface{}) []byte {
	switch v := argValue.(type) {
	case *ast.Variable:
		variableName := v.GetValue().(*ast.Name).Value
		return coerceVariableToBuffer(variables[variableName])

	case *ast.ListValue:
		list := v.GetValues().([]ast.Value)
		buf := new(bytes.Buffer)
		buf.WriteByte('[')
		for _, item := range list {
			itemValue := toCorrectGraphQLArgument(item, variables)
			buf.Write(itemValue)
			buf.WriteByte(',')
		}
		buf.WriteByte(']')

		return buf.Bytes()
	case *ast.IntValue:
		return []byte(fmt.Sprintf("%s", v.GetValue()))
	case *ast.FloatValue:
		return []byte(fmt.Sprintf("%f", v.GetValue()))
	case *ast.StringValue:
		return []byte(strconv.Quote(fmt.Sprintf("%s", v.GetValue())))
	case *ast.EnumValue:
		return []byte(fmt.Sprintf("%s", v.GetValue()))
	case *ast.BooleanValue:
		return []byte(fmt.Sprintf("%s", v.GetValue()))
	case *ast.ObjectValue:
		panic("not supported yet")
	}

	panic("failed to assert graphql argument")
}

func coerceVariableToBuffer(value interface{}) []byte {
	switch v := value.(type) {
	case int:
		return []byte(fmt.Sprintf("%d", v))
	case string:
		return []byte(strconv.Quote(fmt.Sprintf("%s", v)))
	case bool:
		return []byte(fmt.Sprintf("%v", v))
	case []interface{}:
		buf := new(bytes.Buffer)
		buf.WriteByte('[')
		for _, vv := range v {
			buf.Write(coerceVariableToBuffer(vv))
			buf.WriteByte(',')
		}
		buf.WriteByte(']')

		return buf.Bytes()
	default:
		panic("unknown type")
	}
}
