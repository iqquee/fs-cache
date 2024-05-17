package fscache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// memdis test cases
var memdisTestCases = []map[string]MemdisData{
	{
		"key1": MemdisData{
			Value:    "value1",
			Duration: time.Now().Add(time.Minute),
		},
	},
	{
		"key2": MemdisData{
			Value:    10,
			Duration: time.Time{},
		},
	},
	{
		"key3": MemdisData{
			Value:    true,
			Duration: time.Time{},
		},
	},
}

func TestSet(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	err := ch.Memdis().Set("missing_key", "value1", time.Minute)
	require.NoError(t, err)

	err = ch.Memdis().Set("key1", "value1", time.Minute)
	require.Error(t, err) // this one exists, it raises an error
}

func TestGet(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	value, err := ch.Memdis().Get("key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", value)

	_, err = ch.Memdis().Get("missing_key")
	require.ErrorIs(t, err, ErrKeyNotFound)
}

func TestDel(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	err := ch.Memdis().Del("key1")
	require.NoError(t, err)
}

func TestClear(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	err := ch.Memdis().Clear()
	require.NoError(t, err)
}

func TestSize(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	value := ch.Memdis().Size()
	assert.EqualValues(t, 3, value)
}

func TestDebug(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	ch.Debug()
	assert.True(t, debug)
}

func TestOverWrite(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	if err := ch.Memdis().OverWrite("key1", "overwrite1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestOverWriteWithKey(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	if err := ch.Memdis().OverWriteWithKey("key1", "newKey1", "value1", time.Minute); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}

func TestTypeOf(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	typeOf, err := ch.Memdis().TypeOf("key1")
	if err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, typeOf)
}

func TestKeyValuePairs(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	data := ch.Memdis().KeyValuePairs()
	assert.NotNil(t, data)
}

func TestSetMany(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	testCase := []map[string]MemdisData{
		{
			"key4": MemdisData{
				Value:    "value4",
				Duration: time.Now().Add(time.Minute),
			},
			"key5": MemdisData{
				Value: false,
			},
		},
	}

	data, err := ch.Memdis().SetMany(testCase)
	require.NoError(t, err)
	assert.NotNil(t, data)
}

func TestGetMany(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	keys := []string{"key1", "key2"}

	result := ch.Memdis().GetMany(keys)
	assert.NotNil(t, result)
}

func TestKeys(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	keys := ch.Memdis().Keys()
	assert.NotNil(t, keys)
}

func TestValues(t *testing.T) {
	md := Memdis{
		storage: memdisTestCases,
	}
	ch := Cache{
		MemdisInstance: md,
	}

	values := ch.Memdis().Values()
	assert.NotNil(t, values)
}
