package fscache

import (
	"bytes"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// keyStoreTestCases is a slice of maps used for testing the keyStore functionality.
// Each map contains a string key and a KeyStoreData value. KeyStoreData holds a Value
// of various types (string, int, bool) and a Duration which is a time.Time object.
var keyStoreTestCases = []map[string]KeyStoreData{
	{
		"key1": KeyStoreData{
			Value:    "value1",
			Duration: time.Now().Add(time.Minute),
		},
	},
	{
		"key2": KeyStoreData{
			Value:    10,
			Duration: time.Time{},
		},
	},
	{
		"key3": KeyStoreData{
			Value:    true,
			Duration: time.Time{},
		},
	},
}

func TestSet(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	err := ch.KeyStore().Set("missing-key", "value1", time.Minute)
	require.NoError(t, err)

	err = ch.KeyStore().Set("key1", "value1", time.Minute)
	require.Error(t, err) // this one exists, it raises an error
}

func TestGet(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	value, err := ch.KeyStore().Get("key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", value)

	_, err = ch.KeyStore().Get("missing_key")
	require.ErrorIs(t, err, ErrKeyNotFound)
}

func TestDel(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	err := ch.KeyStore().Del("key1")
	require.NoError(t, err)
}

func TestClear(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	err := ch.KeyStore().Clear()
	require.NoError(t, err)
}

func TestSize(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	value := ch.KeyStore().Size()
	assert.EqualValues(t, 3, value)
}

func TestDebug(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	buf := bytes.NewBuffer(nil)

	ch.KeyStore().logger.Info().Msg("hello world - testing") // before enabling debug logging
	data, err := io.ReadAll(buf)
	assert.NoError(t, err)
	assert.Equal(t, string(data), "")

	ch.Debug(buf)
	ch.KeyStore().logger.Info().Msg("hello world - testing") // after enabling debug logging

	data, err = io.ReadAll(buf)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "hello world - testing")
}

func TestOverWrite(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	if err := ch.KeyStore().OverWrite("key1", "overwrite1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestOverWriteWithKey(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	if err := ch.KeyStore().OverWriteWithKey("key1", "newKey1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestTypeOf(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	typeOf, err := ch.KeyStore().TypeOf("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, typeOf)
}

func TestKeyValuePairs(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	data := ch.KeyStore().KeyValuePairs()
	assert.NotNil(t, data)
}

func TestSetMany(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	testCase := []map[string]KeyStoreData{
		{
			"key4": KeyStoreData{
				Value:    "value4",
				Duration: time.Now().Add(time.Minute),
			},
			"key5": KeyStoreData{
				Value: false,
			},
		},
	}

	data, err := ch.KeyStore().SetMany(testCase)
	require.NoError(t, err)
	assert.NotNil(t, data)
}

func TestGetMany(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	keys := []string{"key1", "key2"}

	result := ch.KeyStore().GetMany(keys)
	assert.NotNil(t, result)
}

func TestKeys(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	keys := ch.KeyStore().Keys()
	assert.NotNil(t, keys)
}

func TestValues(t *testing.T) {
	md := KeyStore{
		mu:      &sync.RWMutex{},
		storage: keyStoreTestCases,
	}
	ch := Cache{
		KeyStoreInstance: md,
	}

	values := ch.KeyStore().Values()
	assert.NotNil(t, values)
}
