package fscache

import (
	"errors"
	"reflect"
	"time"
)

var (
	// ErrKeyNotFound key not found
	ErrKeyNotFound = errors.New("key not found")
	// ErrKeyExists key already exists
	ErrKeyExists = errors.New("key already exist")
)

// Set() adds a new data into the in-memory storage
func (md *Memdis) Set(key string, value any, duration ...time.Duration) error {
	md.mu.Lock()
	defer md.mu.Unlock()

	for _, cache := range md.storage {
		if _, ok := cache[key]; ok {
			return ErrKeyExists
		}
	}

	var ttl time.Duration
	for i, v := range duration {
		if i == 0 {
			ttl = v
			break
		}
	}

	fs := make(map[string]MemdisData)
	fs[key] = MemdisData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	md.storage = append(md.storage, fs)

	return nil
}

// SetMany() sets many data objects into memory for later access
func (md *Memdis) SetMany(data []map[string]MemdisData) ([]map[string]any, error) {
	md.mu.Lock()
	defer md.mu.Unlock()

	md.storage = append(md.storage, data...)
	KeyValuePairs := md.KeyValuePairs()

	return KeyValuePairs, nil
}

// Get() retrieves a data from the in-memory storage
func (md *Memdis) Get(key string) (any, error) {
	md.mu.Lock()
	defer md.mu.Unlock()

	for _, cache := range md.storage {
		if val, ok := cache[key]; ok {
			return val.Value, nil
		}
	}

	return nil, ErrKeyNotFound
}

// GetMany() retrieves data with matching keys from the in-memory storage
func (md *Memdis) GetMany(keys []string) []map[string]any {
	keyValuePairs := []map[string]any{}
	for _, cache := range md.storage {
		data := make(map[string]any)
		for _, key := range keys {
			if val, ok := cache[key]; ok {
				data[key] = val.Value
				keyValuePairs = append(keyValuePairs, data)
			}
		}
	}

	return keyValuePairs
}

// Del() deletes a data from the in-memory storage
func (md *Memdis) Del(key string) error {
	md.mu.Lock()
	defer md.mu.Unlock()

	for index, cache := range md.storage {
		if _, ok := cache[key]; ok {
			md.storage = append(md.storage[:index], md.storage[index+1:]...)
			return nil
		}
	}

	return ErrKeyNotFound
}

// Clear() deletes all data from the in-memory storage
func (md *Memdis) Clear() error {
	md.storage = md.storage[:0]
	return nil
}

// Size() retrieves the total data objects in the in-memory storage
func (md *Memdis) Size() int {
	return len(md.storage)
}

// OverWrite() updates an already set value using it key
func (md *Memdis) OverWrite(key string, value any, duration ...time.Duration) error {
	md.mu.Lock()
	defer md.mu.Unlock()

	var isFound bool
	for index, cache := range md.storage {
		if _, ok := cache[key]; ok {
			isFound = true
			md.storage = append(md.storage[:index], md.storage[index+1:]...)
		}
	}

	if !isFound {
		return ErrKeyNotFound
	}

	var ttl time.Duration
	for i, v := range duration {
		if i == 0 {
			ttl = v
			break
		}
	}

	fs := make(map[string]MemdisData)
	fs[key] = MemdisData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	md.storage = append(md.storage, fs)

	return nil
}

// OverWriteWithKey() updates an already set value and key using the previously set key
func (md *Memdis) OverWriteWithKey(prevkey, newKey string, value any, duration ...time.Duration) error {
	md.mu.Lock()
	defer md.mu.Unlock()

	var isFound bool
	for index, cache := range md.storage {
		if _, ok := cache[prevkey]; ok {
			isFound = true
			md.storage = append(md.storage[:index], md.storage[index+1:]...)
		}
	}

	if !isFound {
		return ErrKeyNotFound
	}

	var ttl time.Duration
	for i, v := range duration {
		if i == 0 {
			ttl = v
			break
		}
	}

	fs := make(map[string]MemdisData)
	fs[newKey] = MemdisData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	md.storage = append(md.storage, fs)

	return nil
}

// Keys() returns all the keys in the storage
func (md *Memdis) Keys() []string {
	var keys []string
	for _, cache := range md.storage {
		for key := range cache {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values() returns all the values in the storage
func (md *Memdis) Values() []any {
	var values []any
	for _, cache := range md.storage {
		for _, v := range cache {
			values = append(values, v.Value)
		}
	}

	return values
}

// TypeOf() returns the data type of a value
func (md *Memdis) TypeOf(key string) (string, error) {
	for _, cache := range md.storage {
		value, ok := cache[key]
		if ok {
			return reflect.TypeOf(value.Value).String(), nil
		}
	}

	return "", ErrKeyNotFound
}

// KeyValuePairs() returns an array of key value pairs of all the data in the storage
func (md *Memdis) KeyValuePairs() []map[string]any {
	keyValuePairs := []map[string]any{}

	for _, v := range md.storage {
		data := make(map[string]any)
		for key, value := range v {
			data[key] = value.Value
		}

		keyValuePairs = append(keyValuePairs, data)
	}

	return keyValuePairs
}
