package fscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// test cases
var testCases = []map[string]CacheData{
	{
		"key1": CacheData{
			Value:    "value1",
			Duration: time.Now().Add(time.Minute),
		},
	},
	{
		"key2": CacheData{
			Value:    10,
			Duration: time.Time{},
		},
	},
	{
		"key3": CacheData{
			Value:    true,
			Duration: time.Time{},
		},
	},
}

func TestSet(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	if err := ch.KeyValuePair().Set("key1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestGet(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	value, err := ch.KeyValuePair().Get("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.EqualValues(t, "value1", value)
}

func TestDel(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	if err := ch.KeyValuePair().Del("key1"); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestClear(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	if err := ch.KeyValuePair().Clear(); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestSize(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	value := ch.KeyValuePair().Size()
	assert.EqualValues(t, 3, value)
}

func TestDebug(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	ch.Debug()
	assert.EqualValues(t, true, debug)
}

func TestOverWrite(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	if err := ch.KeyValuePair().OverWrite("key1", "overwrite1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestOverWriteWithKey(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	if err := ch.KeyValuePair().OverWriteWithKey("key1", "newKey1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestTypeOf(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	typeOf, err := ch.KeyValuePair().TypeOf("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, typeOf)
}

func TestKeyValuePairs(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	datas := ch.KeyValuePair().KeyValuePairs()
	assert.NotNil(t, datas)
}

func TestSetMany(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	testCase := []map[string]CacheData{
		{
			"key4": CacheData{
				Value:    "value4",
				Duration: time.Now().Add(time.Minute),
			},
			"key5": CacheData{
				Value: false,
			},
		},
	}

	datas, err := ch.KeyValuePair().SetMany(testCase)
	if err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, datas)
}

func TestGetMany(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	keys := []string{"key1", "key2"}

	result := ch.KeyValuePair().GetMany(keys)
	assert.NotNil(t, result)
}

func TestKeys(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	keys := ch.KeyValuePair().Keys()
	assert.NotNil(t, keys)
}

func TestValues(t *testing.T) {
	kp := KeyPair{
		Storage: testCases,
	}
	ch := Cache{
		KeyPair: kp,
	}

	values := ch.KeyValuePair().Values()
	assert.NotNil(t, values)
}
