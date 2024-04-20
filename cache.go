package fscache

import (
	"fmt"
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
	ch := Cache{Fscache: fs}

	// run go routine to check if the duration has expired and then delete it from off the array
	go func() {
		tt := time.Now()
		for i, v := range ch.Fscache {
			if debug {
				fmt.Println("go routine running...")
			}

			cache := v["duration"]
			if tt.Before(cache.duration) {
				if err := ch.delIndex(i); err != nil {
					if debug {
						fmt.Printf("[error deleting after Expire ::: %v]", err)
					}
				}
			}
		}
	}()

	op := Operations(&ch)

	return op
}

// delIndex is used internally to delete a set object by its index
func (ch *Cache) delIndex(index int) error {
	for i := range ch.Fscache {
		if index == i {
			ch.Fscache = append(ch.Fscache[:index], ch.Fscache[index+1:]...)
			return nil
		}
	}

	return errKeyNotFound
}
