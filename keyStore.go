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
func (ks *KeyStore) Set(key string, value any, duration ...time.Duration) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	for _, cache := range ks.storage {
		if _, ok := cache[key]; ok {
			return ErrKeyExists
		}
	}

	var ttl time.Duration
	if len(duration) > 0 {
		ttl = duration[0]
	}

	fs := make(map[string]KeyStoreData)
	fs[key] = KeyStoreData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	ks.storage = append(ks.storage, fs)

	return nil
}

// SetMany() sets many data objects into memory for later access
func (ks *KeyStore) SetMany(data []map[string]KeyStoreData) ([]map[string]any, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.storage = append(ks.storage, data...)
	KeyValuePairs := ks.KeyValuePairs()

	return KeyValuePairs, nil
}

// Get() retrieves a data from the in-memory storage
func (ks *KeyStore) Get(key string) (any, error) {
	for _, cache := range ks.storage {
		if val, ok := cache[key]; ok {
			return val.Value, nil
		}
	}

	return nil, ErrKeyNotFound
}

// GetMany() retrieves data with matching keys from the in-memory storage
func (ks *KeyStore) GetMany(keys []string) []map[string]any {
	keyValuePairs := []map[string]any{}
	for _, cache := range ks.storage {
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
func (ks *KeyStore) Del(key string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	for index, cache := range ks.storage {
		if _, ok := cache[key]; ok {
			ks.storage = append(ks.storage[:index], ks.storage[index+1:]...)
			return nil
		}
	}

	return ErrKeyNotFound
}

// Clear() deletes all data from the in-memory storage
func (ks *KeyStore) Clear() error {
	ks.storage = ks.storage[:0]
	return nil
}

// Size() retrieves the total data objects in the in-memory storage
func (ks *KeyStore) Size() int {
	return len(ks.storage)
}

// OverWrite() updates an already set value using it key
func (ks *KeyStore) OverWrite(key string, value any, duration ...time.Duration) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var isFound bool
	for index, cache := range ks.storage {
		if _, ok := cache[key]; ok {
			isFound = true
			ks.storage = append(ks.storage[:index], ks.storage[index+1:]...)
		}
	}

	if !isFound {
		return ErrKeyNotFound
	}

	var ttl time.Duration
	if len(duration) > 0 {
		ttl = duration[0]
	}

	fs := make(map[string]KeyStoreData)
	fs[key] = KeyStoreData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	ks.storage = append(ks.storage, fs)

	return nil
}

// OverWriteWithKey() updates an already set value and key using the previously set key
func (ks *KeyStore) OverWriteWithKey(prevkey, newKey string, value any, duration ...time.Duration) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var isFound bool
	for index, cache := range ks.storage {
		if _, ok := cache[prevkey]; ok {
			isFound = true
			ks.storage = append(ks.storage[:index], ks.storage[index+1:]...)
		}
	}

	if !isFound {
		return ErrKeyNotFound
	}

	var ttl time.Duration
	if len(duration) > 0 {
		ttl = duration[0]
	}

	fs := make(map[string]KeyStoreData)
	fs[newKey] = KeyStoreData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	ks.storage = append(ks.storage, fs)

	return nil
}

// Keys() returns all the keys in the storage
func (ks *KeyStore) Keys() []string {
	var keys []string
	for _, cache := range ks.storage {
		for key := range cache {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values() returns all the values in the storage
func (ks *KeyStore) Values() []any {
	var values []any
	for _, cache := range ks.storage {
		for _, v := range cache {
			values = append(values, v.Value)
		}
	}

	return values
}

// TypeOf() returns the data type of a value
func (ks *KeyStore) TypeOf(key string) (string, error) {
	for _, cache := range ks.storage {
		value, ok := cache[key]
		if ok {
			return reflect.TypeOf(value.Value).String(), nil
		}
	}

	return "", ErrKeyNotFound
}

// KeyValuePairs() returns an array of key value pairs of all the data in the storage
func (ks *KeyStore) KeyValuePairs() []map[string]any {
	keyValuePairs := []map[string]any{}

	for _, v := range ks.storage {
		data := make(map[string]any)
		for key, value := range v {
			data[key] = value.Value
		}

		keyValuePairs = append(keyValuePairs, data)
	}

	return keyValuePairs
}
