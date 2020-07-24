package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/terra-project/mantle/types"
)

func GenerateQuery(query types.GraphQLQueryRaw) string {
	var buf bytes.Buffer
	buf.Write([]byte("query"))
	generate(&buf, reflect.TypeOf(query), false)

	return buf.String()
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

func generate(w io.Writer, t reflect.Type, inline bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		generate(w, t.Elem(), false)
	case reflect.Struct:
		// If the type implements json.Unmarshaler, it's a scalar. Don't expand it.
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			io.WriteString(w, "{")
		}
		for i := 0; i < t.NumField(); i++ {
			if i != 0 {
				io.WriteString(w, ",")
			}
			f := t.Field(i)
			value, ok := f.Tag.Lookup("mantle")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					// if there is a character before (,
					// then it's an alias
					if i := strings.Index(value, "("); i != 0 {
						io.WriteString(w, fmt.Sprintf("%s:%s", f.Name, value))
					} else {
						io.WriteString(w, fmt.Sprintf("%s%s", f.Name, value))
					}

				} else {
					io.WriteString(w, f.Name)
				}
			}
			generate(w, f.Type, inlineField)
		}
		if !inline {
			io.WriteString(w, "}")
		}
	}
}
