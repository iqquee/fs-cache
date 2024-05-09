package fscache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var noSqlTestCases = []interface{}{
	struct {
		Name string
		Age  int
	}{
		Name: "Jane Doe",
		Age:  25,
	},
	map[string]interface{}{
		"name":    "John Doe",
		"age":     35,
		"colName": "users",
	},
	map[string]interface{}{
		"name":    "Jane Dice",
		"age":     35,
		"colName": "users",
	},
}

func Test_Collection(t *testing.T) {
	ch := Cache{}

	col := ch.NoSql().Collection("user")
	assert.NotNil(t, col)
	assert.Equal(t, "users", col.collectionName)
}

func Test_Insert(t *testing.T) {
	ch := Cache{}

	var counter int
	name := fmt.Sprintf("testCase_%v", counter+1)
	for _, v := range noSqlTestCases {
		t.Run(name, func(t *testing.T) {
			res, err := ch.NoSql().Collection("user").Insert(v)
			if err != nil {
				assert.Error(t, err)
			}

			assert.NotNil(t, v, res)
		})

		counter++
	}
}

func Test_InsertMany(t *testing.T) {
	ch := Cache{}

	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
}

func Test__Filter_First(t *testing.T) {
	ch := Cache{}

	// insert a new record
	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)

	filters := []map[string]interface{}{
		{"age": 35.0}, // filter out records of age 35

		nil, // for nil params
	}

	for _, v := range filters {
		var name string
		if v == nil {
			name = "nil params"
		} else {
			name = "not nil params"
		}

		t.Run(name, func(t *testing.T) {
			result, err := ch.NoSql().Collection("users").Filter(v).First()
			if err != nil {
				assert.Error(t, err)
			}

			if result == nil {
				assert.Equal(t, errors.New("filter params cannot be nil"), err)
			}
		})
	}
}

func Test_Filter_All(t *testing.T) {
	ch := Cache{}

	// insert a new record
	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)

	filters := []map[string]interface{}{
		{"age": 35.0}, // filter out records of age 35

		nil, // for nil params
	}

	for _, v := range filters {
		var name string
		if v == nil {
			name = "nil params"
		} else {
			name = "not nil params"
		}

		t.Run(name, func(t *testing.T) {
			result, err := ch.NoSql().Collection("users").Filter(v).All()
			if err != nil {
				assert.Error(t, err)
			}

			assert.NotNil(t, result)
		})
	}
}

func Test_Delete_One(t *testing.T) {
	ch := Cache{}

	// insert a new record
	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)

	filters := []map[string]interface{}{
		{"age": 35.0}, // filter out record of age 35

		nil, // for nil params
	}

	for _, v := range filters {
		var name string
		if v == nil {
			name = "nil params"
		} else {
			name = "not nil params"
		}

		t.Run(name, func(t *testing.T) {
			err := ch.NoSql().Collection("users").Delete(v).One()
			if err != nil {
				assert.Error(t, err)
			}

		})
	}
}

func Test_Delete_All(t *testing.T) {
	ch := Cache{}

	// insert a new record
	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NoError(t, err)

	filters := []map[string]interface{}{
		{"age": 35.0}, // filter out records of age 35
		nil,           // for nil params
	}

	for _, v := range filters {
		var name string
		if v == nil {
			name = "nil params"
		} else {
			name = "not nil params"
		}

		t.Run(name, func(t *testing.T) {
			err := ch.NoSql().Collection("users").Delete(v).All()
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}
