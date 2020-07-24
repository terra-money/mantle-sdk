package types

type DatabaseDriver interface {
	Get(entity interface{}, args map[string]interface{}) interface{}
	Set(entity interface{}) interface{}
	Delete(entity interface{})
}