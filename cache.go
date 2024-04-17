package fscache

import "time"

type (
	// Store is the in-memory storage cache
	// Fscache []map[string]interface{}
	Debug bool

	Cache struct {
		Fscache []map[string]interface{}
	}

	// Operations lists all available operations on the fscache
	Operations interface {
		Set(key string, value interface{}, duration ...time.Time) error
		Get(key string) (interface{}, error)
		Del(key string) error
		Clear() error
		Size() int
		MemSize() int
	}
)

// New initializes an instance of the in-memory storage cache
func New() *Operations {
	var fs []map[string]interface{}

	op := Operations(&Cache{Fscache: fs})
	return &op
}
