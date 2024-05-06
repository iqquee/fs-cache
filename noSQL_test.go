package fscache

import (
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
		"name": "john Doe",
		"age":  35,
	},
}

func Test_Collection(t *testing.T) {
	ns := NoSQL{
		storage: noSqlTestCases,
	}

	ch := Cache{
		NoSQL: ns,
	}

	col := ch.NoSql().Collection("user")
	assert.NotNil(t, col)
	assert.Equal(t, "users", col.collectionName)
}

func Test_Insert(t *testing.T) {
	ns := NoSQL{
		storage: []interface{}{},
	}

	ch := Cache{
		NoSQL: ns,
	}

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
	ns := NoSQL{
		storage: []interface{}{},
	}

	ch := Cache{
		NoSQL: ns,
	}

	err := ch.NoSql().Collection("user").InsertMany(noSqlTestCases)
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
}
