package fscache

import (
	"time"
)

type (
	// cacheData object
	cacheData struct {
		value    interface{}
		duration time.Time
	}

	// Cache object instance
	Cache struct {
		// debug enables debugging
		debug   bool
		Fscache []map[string]cacheData
	}

	// Operations lists all available operations on the fscache
	Operations interface {
		// Set() adds a new data into the in-memmory storage
		Set(key string, value interface{}, duration ...time.Duration) error
		// Get() retrieves a data from the in-memmory storage
		Get(key string) (interface{}, error)
		// Del() deletes a data from the in-memmory storage
		Del(key string) error
		// Clear() deletes all datas from the in-memmory storage
		Clear() error
		// Size() retrieves the total data objects in the in-memmory storage
		Size() int
		// Debug() enables debug to get certain logs
		Debug()
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	var fs []map[string]cacheData
	ch := Cache{
		Fscache: fs,
	}

	op := Operations(&ch)
	return op
}
