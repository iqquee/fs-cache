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
		Set(key string, value interface{}, duration ...time.Duration) error
		Get(key string) (interface{}, error)
		Del(key string) error
		Clear() error
		Size() int
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
			cache := v["duration"]
			if tt.Before(cache.duration) {
				if err := ch.delIndex(i); err != nil {
					fmt.Printf("[error deleting after Expire ::: %v]", err)
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
