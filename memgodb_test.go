package fscache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var MemgodbTestCases = []any{
	struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Jane Doe",
		Age:  25,
	},
	map[string]any{
		"name":    "John Doe",
		"age":     35,
		"colName": "users",
	},
	map[string]any{
		"name":    "Jane Dice",
		"age":     35,
		"colName": "users",
	},
}

func Test_Collection(t *testing.T) {
	ch := Cache{}

	col := ch.Memgodb().Collection("user")
	assert.NotNil(t, col)
	assert.Equal(t, "users", col.collectionName)
}

func Test_Insert_One(t *testing.T) {
	ch := Cache{}

	var counter int
	name := fmt.Sprintf("testCase_%v", counter+1)
	for _, v := range MemgodbTestCases {
		t.Run(name, func(t *testing.T) {
			res, err := ch.Memgodb().Collection("user").Insert(v).One()
			if err != nil {
				assert.Error(t, err)
			}

			assert.NotNil(t, v, res)
		})

		counter++
	}
}

func Test_Insert_Many(t *testing.T) {
	ch := Cache{}

	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, res)
}

func Test_Insert_FromJsonFile(t *testing.T) {
	ch := Cache{}

	testCases := []struct {
		fileName      string
		expectedError error
		name          string
		message       string
	}{
		{
			fileName:      "./testJsonFiles/objects.json",
			name:          "objects [slice] file",
			expectedError: nil,
			message:       "success",
		},
		{
			fileName:      "./testJsonFiles/object.json",
			name:          "object [map] file",
			expectedError: nil,
			message:       "success",
		},
		{
			fileName:      "./testJsonFiles/string.json",
			name:          "[string] file",
			expectedError: errors.New("file must contain either an array of [objects ::: slice] or [object ::: map]"),
			message:       "fail",
		},
		{
			fileName:      "./testJsonFiles/empty.json",
			name:          "[empty] file",
			expectedError: errors.New("invalid json file"),
			message:       "fail",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ch.Memgodb().Collection("user").Insert(nil).FromJsonFile(testCase.fileName)
			if testCase.message != "success" {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}
		})
	}
}

func Test__Filter_First(t *testing.T) {
	ch := Cache{}

	// insert a new records
	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, res)

	testCases := []struct {
		name          string
		expectedError error
		message       string
		filter        map[string]any
	}{
		{
			name:          "nil params",
			expectedError: errors.New("filter params cannot be nil"),
			message:       "fail_1",
			filter:        nil, // for nil params
		},
		{
			name:          "not nil params",
			expectedError: nil,
			message:       "success",
			filter:        map[string]any{"age": 35.0}, // filter out records of age 35
		},
		{
			name:          "incorrect params",
			expectedError: errors.New("record not found"),
			message:       "fail_2",
			filter:        map[string]any{"age": 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := ch.Memgodb().Collection("users").Filter(testCase.filter).First()
			if testCase.message == "fail_1" {
				assert.Equal(t, testCase.expectedError, err)
			} else if testCase.message == "fail_2" {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func Test_Filter_All(t *testing.T) {
	ch := Cache{}

	// insert a new records
	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, res)

	testCases := []struct {
		name          string
		expectedError error
		message       string
		filter        map[string]any
	}{
		{
			name:          "not nil params",
			expectedError: nil,
			message:       "success",
			filter:        map[string]any{"age": 35.0}, // filter out records of age 35
		},
		{
			name:          "incorrect params",
			expectedError: errors.New("record not found"),
			message:       "fail",
			filter:        map[string]any{"age": 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			results, err := ch.Memgodb().Collection("users").Filter(testCase.filter).All()
			if testCase.message != "success" {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NotNil(t, results)
			}
		})
	}
}

func Test_Delete_One(t *testing.T) {
	ch := Cache{}

	// insert a new record
	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, res)

	filters := []map[string]any{
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
			err := ch.Memgodb().Collection("users").Delete(v).One()
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func Test_Delete_All(t *testing.T) {
	ch := Cache{}

	// insert a new record
	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, res)

	filters := []map[string]any{
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
			err := ch.Memgodb().Collection("users").Delete(v).All()
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func Test_Update_One(t *testing.T) {
	ch := Cache{}

	// insert a new record
	res, err := ch.Memgodb().Collection("user").Insert(nil).Many(MemgodbTestCases)
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, res)

	testCases := []struct {
		expectedError error
		name          string
		filter        map[string]any
		message       string
		update        map[string]any
	}{
		{
			name:          "correct filter params",
			expectedError: nil,
			message:       "success",
			filter: map[string]any{
				"age": 35.0,
			},
			update: map[string]any{
				"age": 29,
			},
		},
		{
			name:          "nil filter params",
			expectedError: errors.New("filter params cannot be nil"),
			message:       "failed_1",
			filter:        nil,
			update: map[string]any{
				"age": 28,
			},
		},
		{
			name:          "not found params",
			expectedError: errors.New("record not found"),
			message:       "failed_2",
			filter: map[string]any{
				"age": 300.0,
			},
			update: map[string]any{
				"age": 28,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ch.Memgodb().Collection("user").Update(testCase.filter, testCase.update).One()
			if testCase.message == "success" {
				assert.Equal(t, testCase.expectedError, err)
			} else if testCase.message == "fail_1" {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.Equal(t, testCase.expectedError, err)
			}
		})
	}
}

func Test_Persist(t *testing.T) {
	ch := Cache{}

	err := ch.Memgodb().Persist()
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
}

func Test_LoadDefault(t *testing.T) {
	ch := Cache{}

	err := ch.Memgodb().LoadDefault()
	if err != nil {
		assert.Equal(t, errors.New("error finding file"), err)
	}
}
