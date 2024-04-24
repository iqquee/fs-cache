package fscache

import (
	"errors"
	"reflect"
	"time"
)

var (
	// errKeyNotFound key not found
	errKeyNotFound = errors.New("key not found")
	// errKeyExists key already exists
	errKeyExists = errors.New("key already exist")
)

// Set() adds a new data into the in-memmory storage
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

// Get() retrieves a data from the in-memmory storage
func (ch *Cache) Get(key string) (interface{}, error) {
	for _, cache := range ch.Fscache {
		if val, ok := cache[key]; ok {
			return val.value, nil
		}
	}

	return nil, errKeyNotFound
}

// Del() deletes a data from the in-memmory storage
func (ch *Cache) Del(key string) error {
	var isFound bool
	for index, cache := range ch.Fscache {
		if _, ok := cache[key]; ok {
			isFound = true
			ch.Fscache = append(ch.Fscache[:index], ch.Fscache[index+1:]...)
			return nil
		}
	}

	if !isFound {
		return errKeyNotFound
	}

	return errKeyNotFound
}

// Clear() deletes all datas from the in-memmory storage
func (ch *Cache) Clear() error {
	ch.Fscache = ch.Fscache[:0]

	return nil
}

// Size() retrieves the total data objects in the in-memmory storage
func (ch *Cache) Size() int {
	return len(ch.Fscache)
}

// Debug() enables debug to get certain logs
func (ch *Cache) Debug() {
	ch.debug = true
}

// OverWrite() updates an already set value using it key
func (ch *Cache) OverWrite(key string, value interface{}, duration ...time.Duration) error {
	var isFound bool
	for index, cache := range ch.Fscache {
		if _, ok := cache[key]; ok {
			isFound = true
			ch.Fscache = append(ch.Fscache[:index], ch.Fscache[index+1:]...)
		}
	}

	if !isFound {
		return errKeyNotFound
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

// OverWriteWithKey() updates an already set value and key using the previously set key
func (ch *Cache) OverWriteWithKey(prevkey, newKey string, value interface{}, duration ...time.Duration) error {
	var isFound bool
	for index, cache := range ch.Fscache {
		if _, ok := cache[prevkey]; ok {
			isFound = true
			ch.Fscache = append(ch.Fscache[:index], ch.Fscache[index+1:]...)
		}
	}

	if !isFound {
		return errKeyNotFound
	}

	var ttl time.Duration
	for i, v := range duration {
		if i == 0 {
			ttl = v
			break
		}
	}

	fs := make(map[string]cacheData)
	fs[newKey] = cacheData{
		value:    value,
		duration: time.Now().Add(ttl),
	}

	ch.Fscache = append(ch.Fscache, fs)

	return nil
}

// ExportJson() exports all saves data objects as json
func (ch *Cache) ExportJson() []map[string]cacheData {
	return nil
}

// ImportJson() takes in an array of json objects and saves it into memory for later access
func (ch *Cache) ImportJson([]map[string]interface{}) error {
	return nil
}

// Keys() returns all the keys in the storage
func (ch *Cache) Keys() []string {
	var keys []string
	for i := 0; i < len(ch.Fscache); i++ {
		for key := range ch.Fscache[i] {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values() returns all the values in the storage
func (ch *Cache) Values() []interface{} {
	var values []interface{}
	for i := 0; i < len(ch.Fscache); i++ {
		for _, v := range ch.Fscache[i] {
			values = append(values, v.value)
		}
	}

	return values
}

// TypeOf() returns the data type of a value
func (ch *Cache) TypeOf(key string) (string, error) {
	for _, cache := range ch.Fscache {
		value, ok := cache[key]
		if ok {
			return reflect.TypeOf(value.value).String(), nil
		}
	}

	return "", errKeyNotFound
}

// SaveToFile() saves the array of objects into a file
func (ch *Cache) SaveToFile(fileName string) error {
	return nil
}
