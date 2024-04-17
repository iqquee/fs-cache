package fscache

import (
	"errors"
	"time"
)

var (
	errKeyNotFound = errors.New("key was not found")
)

// Set()
func (ch *Cache) Set(key string, value interface{}, duration ...time.Time) error {
	fs := make(map[string]interface{})
	fs[key] = value
	ch.Fscache = append(ch.Fscache, fs)

	return nil
}

// Get()
func (ch *Cache) Get(key string) (interface{}, error) {
	for _, cache := range ch.Fscache {
		if val, ok := cache[key]; ok {
			return val, nil
		}
	}

	return "", errKeyNotFound
}

func (ch *Cache) Del(key string) error {
	for index, cache := range ch.Fscache {
		if _, ok := cache[key]; ok {
			ch.Fscache = append(ch.Fscache[:index], ch.Fscache[index+1:]...)
			return nil
		}
	}

	return errKeyNotFound
}
func (ch *Cache) Clear() error {
	ch.Fscache = ch.Fscache[:0]

	return nil
}

func (ch Cache) Size() int {
	return len(ch.Fscache)
}

func (ch Cache) MemSize() int
