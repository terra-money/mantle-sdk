package constants

var DefaultLimit = 50
var DefaultOrder = false

const (
	ASC = iota
	DESC
)

func GetOrder(args map[string]interface{}) int {
	order, orderExists := args["Order"]
	if !orderExists {
		return ASC
	}

	switch order {
	case "DESC": return DESC
	case "ASC": return ASC
	default: return ASC
	}
}
