package keyvalue

type KeyValue interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Clear() error
}
