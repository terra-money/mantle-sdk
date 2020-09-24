package generate

const (
	DESC = iota
	ASC
)

var defaultLimit = 50
var defaultOrder = DESC

var stringOrderToConst = map[string]int{
	"desc": DESC,
	"DESC": DESC,
	"asc":  ASC,
	"ASC":  ASC,
}
