package generate

import (
	"github.com/terra-project/mantle/utils"
)

var ReservedArgKeys = []string{
	"limit",
	"order",
	"order_by",
	"offset",
}

func FilterArgs(args map[string]interface{}, skip []string) map[string]interface{} {
	next := make(map[string]interface{})
	for argKey, argValue := range args {
		if utils.SliceContainsString(skip, argKey) {
			continue
		}

		next[argKey] = argValue
	}

	return next
}
