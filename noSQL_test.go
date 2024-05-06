package fscache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var noSqlTestCases = []interface{}{
	&struct {
		Name string
		Age  int
	}{
		Name: "Jane Doe",
		Age:  25,
	},
	map[string]interface{}{
		"name": "jane",
		"age":  25,
	},
}

func Test_Collection(t *testing.T) {
	ns := NoSQL{
		storage: noSqlTestCases,
	}

	ch := Cache{
		NoSQL: ns,
	}

	var counter int
	name := fmt.Sprintf("testCase %v", counter+1)
	for _, v := range noSqlTestCases {
		t.Run(name, func(t *testing.T) {
			res, err := ch.NoSql().Collection("user").Insert(v)
			if err != nil {
				assert.Error(t, err)
			}

			fmt.Println(res)
			assert.Equal(t, v, res)
		})

		counter++
	}
}
