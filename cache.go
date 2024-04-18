package fscache

import (
	"time"
)

type (
	// Store is the in-memory storage cache
	// Fscache []map[string]interface{}
	Debug bool

	cacheData struct {
		value    interface{}
		duration time.Duration
	}

	Cache struct {
		Fscache []map[string]cacheData
	}

	// Operations lists all available operations on the fscache
	Operations interface {
		Set(key string, value interface{}, duration ...time.Duration) error
		Get(key string) (interface{}, error)
		Del(key string) error
		Clear() error
		Size() int
		// MemSize() int
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	var fs []map[string]cacheData
	ch := Cache{Fscache: fs}

	go func() {

	}()

	op := Operations(&ch)

	return op
}
