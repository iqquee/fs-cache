package fscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// test cases
var testCases = []map[string]cacheData{
	{
		"key1": cacheData{
			value:    "value1",
			duration: time.Now().Add(time.Minute),
		},
	},
	{
		"key2": cacheData{
			value:    "value2",
			duration: time.Time{},
		},
	},
	{
		"key3": cacheData{
			value:    "value1",
			duration: time.Time{},
		},
	},
}

func TestSet(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	err := ch.Set("key1", "value1", time.Minute)
	if err != nil {
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

	err := ch.Del("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestClear(t *testing.T) {
	ch := Cache{
		Fscache: testCases,
	}

	err := ch.Clear()
	if err != nil {
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
