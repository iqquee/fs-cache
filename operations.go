package fscache

import (
	"errors"
	"time"
)

var (
	errKeyNotFound = errors.New("key not found")
	errKeyExists   = errors.New("key already exist")
)

// Set()
func (ch *Cache) Set(key string, value interface{}, duration ...time.Duration) error {
	for _, cache := range ch.Fscache {
		if _, ok := cache[key]; ok {
			return errKeyExists
		}
	}

	var ttl time.Duration
	for i, v := range duration {
		if i == 0 {
			ttl = v
			break
		}
	}

	fs := make(map[string]cacheData)
	fs[key] = cacheData{
		value:    value,
		duration: time.Now().Add(ttl),
	}

	ch.Fscache = append(ch.Fscache, fs)

	return nil
}

// Get()
func (ch *Cache) Get(key string) (interface{}, error) {
	for _, cache := range ch.Fscache {
		if val, ok := cache[key]; ok {
			return val.value, nil
		}
	}

	return nil, errKeyNotFound
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

func (ch *Cache) Size() int {
	return len(ch.Fscache)
}

// func (ch *Cache) MemSize() int
