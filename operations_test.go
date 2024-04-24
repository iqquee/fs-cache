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
			value:    "value1",
			duration: time.Now().Add(time.Minute),
		},
	},
	{
		"key2": CacheData{
			value:    10,
			duration: time.Time{},
		},
	},
	{
		"key3": CacheData{
			value:    true,
			duration: time.Time{},
		},
	},
}

func TestSet(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	if err := ch.Set("key1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestGet(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	value, err := ch.Get("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.EqualValues(t, "value1", value)
}

func TestDel(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	if err := ch.Del("key1"); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestClear(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	if err := ch.Clear(); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestSize(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	value := ch.Size()
	assert.EqualValues(t, 3, value)
}

func TestDebug(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	ch.Debug()
	assert.EqualValues(t, true, ch.debug)
}

func TestOverWrite(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	if err := ch.OverWrite("key1", "overwrite1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestOverWriteWithKey(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	if err := ch.OverWriteWithKey("key1", "newKey1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestTypeOf(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	typeOf, err := ch.TypeOf("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.EqualValues(t, "string", typeOf)
}

func TestKeyValuePairs(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	datas := ch.KeyValuePairs()
	assert.NotNil(t, datas)
}
