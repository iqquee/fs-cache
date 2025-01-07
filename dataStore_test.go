package fscache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type user struct {
	Name string
	Age  int
}

func TestNameSpace(t *testing.T) {
	fs := New()

	testCases := []struct {
		Name     string
		Params   interface{}
		Expected string
		Schema   Schema
	}{
		{
			Name:     "String Params",
			Params:   "user",
			Expected: "users",
		},
		{
			Name:     "Object Params",
			Params:   user{},
			Expected: "users",
		},
		{
			Name:     "Params value ending with [s]",
			Params:   "atlas", // since it ends with [s], s should not be added to it
			Expected: "atlas",
		},
		{
			Name:     "Params with schema",
			Params:   user{},
			Expected: "users",
			Schema: Schema{
				"Name": "string",
				"Age":  "int",
			},
		},
	}

	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			res := fs.DataStore().Namespace(v.Params, v.Schema)
			assert.Equal(t, v.Expected, res.namespace)
		})
	}
}

func TestCreate(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user", Schema{
		"Name": "string",
		"Age":  "int",
	})

	user := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	err := ns.Create(user)
	assert.NoError(t, err)
}

func TestQuery(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user", Schema{
		"name": "string",
		"age":  "int",
	})

	user := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	err := ns.Create(user)
	assert.NoError(t, err)

	filter := map[string]interface{}{
		"Age": 30,
	}

	res, err := ns.Query(filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestUpdate(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user")

	user := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	err := ns.Create(user)
	assert.NoError(t, err)
	fmt.Println("Created data: ", user)

	filter := map[string]interface{}{
		"Age": 30,
	}

	data := map[string]interface{}{
		"Age": 25,
	}

	err = ns.Update(filter, data)
	assert.NoError(t, err)

	res, err := ns.Query(data)
	assert.NotEmpty(t, res)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user")

	user := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	err := ns.Create(user)
	assert.NoError(t, err)
	fmt.Println("Created data: ", user)

	filter := map[string]interface{}{
		"Age": 30,
	}

	err = ns.Delete(filter)
	assert.NoError(t, err)

	res, err := ns.Query(filter)
	assert.Empty(t, res)
	assert.NoError(t, err)
}

func TestListNamespaces(t *testing.T) {
	fs := New()

	fs.DataStore().Namespace("user")
	fs.DataStore().Namespace("atlas")

	namespaces := fs.DataStore().ListNamespaces()
	assert.NotEmpty(t, namespaces)
}
