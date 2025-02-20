package fscache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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

func TestFirst(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user", Schema{
		"name": "string",
		"age":  "int",
	})

	data1 := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	_ = ns.Create(data1)

	data2 := map[string]interface{}{
		"Name": "John Doe",
		"Age":  35,
	}
	_ = ns.Create(data2)

	filter := map[string]interface{}{
		"Age": 30,
	}

	var response user
	err := ns.First(filter, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Jane Doe", response.Name)
	assert.Equal(t, 30, response.Age)
}

func TestFind(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user", Schema{
		"name": "string",
		"age":  "int",
	})

	data1 := map[string]interface{}{
		"Name": "Jane Doe",
		"Age":  30,
	}
	_ = ns.Create(data1)

	data2 := map[string]interface{}{
		"Name": "John Doe",
		"Age":  30,
	}
	_ = ns.Create(data2)

	filter := map[string]interface{}{
		"Age": 30,
	}

	var response []user
	err := ns.Find(filter, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)
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

func TestConnectSQLDB(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user")

	mockDB := &gorm.DB{}
	ns.ConnectSQLDB(mockDB)
}

func TestConnectMongoDB(t *testing.T) {
	fs := New()

	ns := fs.DataStore().Namespace("user")

	mockDB := &mongo.Database{}
	ns.ConnectMongoDB(mockDB)
}
