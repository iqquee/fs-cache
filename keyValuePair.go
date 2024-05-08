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
func (kp *KeyPair) Set(key string, value interface{}, duration ...time.Duration) error {
	for _, cache := range kp.Storage {
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

	fs := make(map[string]CacheData)
	fs[key] = CacheData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	kp.Storage = append(kp.Storage, fs)

	return nil
}

// SetMany() sets many data objects into memory for later access
func (kp *KeyPair) SetMany(data []map[string]CacheData) ([]map[string]interface{}, error) {
	kp.Storage = append(kp.Storage, data...)
	KeyValuePairs := kp.KeyValuePairs()

	return KeyValuePairs, nil
}

// Get() retrieves a data from the in-memmory storage
func (kp *KeyPair) Get(key string) (interface{}, error) {
	for _, cache := range kp.Storage {
		if val, ok := cache[key]; ok {
			return val.Value, nil
		}
	}

	return nil, errKeyNotFound
}

// GetMany() retrieves datas with matching keys from the in-memmory storage
func (kp *KeyPair) GetMany(keys []string) []map[string]interface{} {
	var keyValuePairs = []map[string]interface{}{}

	for _, cache := range kp.Storage {
		data := make(map[string]interface{})
		for _, key := range keys {
			if val, ok := cache[key]; ok {
				data[key] = val.Value
				keyValuePairs = append(keyValuePairs, data)
			}
		}
	}

	return keyValuePairs
}

// Del() deletes a data from the in-memmory storage
func (kp *KeyPair) Del(key string) error {
	var isFound bool
	for index, cache := range kp.Storage {
		if _, ok := cache[key]; ok {
			isFound = true
			kp.Storage = append(kp.Storage[:index], kp.Storage[index+1:]...)
			return nil
		}
	}

	if !isFound {
		return errKeyNotFound
	}

	return errKeyNotFound
}

// Clear() deletes all datas from the in-memmory storage
func (kp *KeyPair) Clear() error {
	kp.Storage = kp.Storage[:0]

	return nil
}

// Size() retrieves the total data objects in the in-memmory storage
func (kp *KeyPair) Size() int {
	return len(kp.Storage)
}

// OverWrite() updates an already set value using it key
func (kp *KeyPair) OverWrite(key string, value interface{}, duration ...time.Duration) error {
	var isFound bool
	for index, cache := range kp.Storage {
		if _, ok := cache[key]; ok {
			isFound = true
			kp.Storage = append(kp.Storage[:index], kp.Storage[index+1:]...)
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

	fs := make(map[string]CacheData)
	fs[key] = CacheData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	kp.Storage = append(kp.Storage, fs)

	return nil
}

// OverWriteWithKey() updates an already set value and key using the previously set key
func (kp *KeyPair) OverWriteWithKey(prevkey, newKey string, value interface{}, duration ...time.Duration) error {
	var isFound bool
	for index, cache := range kp.Storage {
		if _, ok := cache[prevkey]; ok {
			isFound = true
			kp.Storage = append(kp.Storage[:index], kp.Storage[index+1:]...)
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

	fs := make(map[string]CacheData)
	fs[newKey] = CacheData{
		Value:    value,
		Duration: time.Now().Add(ttl),
	}

	kp.Storage = append(kp.Storage, fs)

	return nil
}

// Keys() returns all the keys in the storage
func (kp *KeyPair) Keys() []string {
	var keys []string
	for _, cache := range kp.Storage {
		for key := range cache {
			keys = append(keys, key)
		}
	}

	return keys
}

// Values() returns all the values in the storage
func (kp *KeyPair) Values() []interface{} {
	var values []interface{}
	for _, cache := range kp.Storage {
		for _, v := range cache {
			values = append(values, v.Value)
		}
	}

	return values
}

// TypeOf() returns the data type of a value
func (kp *KeyPair) TypeOf(key string) (string, error) {
	for _, cache := range kp.Storage {
		value, ok := cache[key]
		if ok {
			return reflect.TypeOf(value.Value).String(), nil
		}
	}

	return "", errKeyNotFound
}

// KeyValuePairs() returns an array of key value pairs of all the datas in the storage
func (kp *KeyPair) KeyValuePairs() []map[string]interface{} {
	var keyValuePairs = []map[string]interface{}{}

	for _, v := range kp.Storage {
		data := make(map[string]interface{})
		for key, value := range v {
			data[key] = value.Value
		}

		keyValuePairs = append(keyValuePairs, data)
	}

	return keyValuePairs
}
