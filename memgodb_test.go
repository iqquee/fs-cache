package fscache

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	fs := New()

	col := fs.Memgodb().Collection("user")
	assert.NotNil(t, col)
	assert.Equal(t, "users", col.collectionName)
}

func Test_Insert_One(t *testing.T) {
	fs := New()

	// insert single records
	var counter int
	name := fmt.Sprintf("testCase_%v", counter+1)
	for _, v := range MemgodbTestCases {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("user").Insert(v)
			require.NoError(t, err)
		})

		counter++
	}

	// to insert many records at once
	err := fs.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)
}

func Test_InsertFromJsonFile(t *testing.T) {
	fs := New()

	testCases := map[string]struct {
		fileName      string
		expectedError error
	}{
		"objects [slice] file": {
			fileName: "./testJsonFiles/objects.json",
		},
		"object [map] file": {
			fileName: "./testJsonFiles/object.json",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("user").InsertFromJsonFile(testCase.fileName)
			assert.Equal(t, testCase.expectedError, err)
		})
	}

	testCases = map[string]struct {
		fileName      string
		expectedError error
	}{
		"[string] file": {
			fileName:      "./testJsonFiles/string.json",
			expectedError: errors.New("file must contain either an array of [objects ::: slice] or [object ::: map]"),
		},
		"[empty] file": {
			fileName:      "./testJsonFiles/empty.json",
			expectedError: errors.New("invalid json file"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("user").InsertFromJsonFile(testCase.fileName)
			require.Error(t, err)
			require.ErrorContains(t, err, testCase.expectedError.Error())
		})
	}
}

func Test__Filter_First(t *testing.T) {
	fs := New()

	// insert a new records
	err := fs.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)

	testCases := map[string]struct {
		expectedError error
		filter        map[string]any
	}{
		"not nil params": {
			filter: map[string]any{"age": 35.0}, // filter out records of age 35
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := fs.Memgodb().Collection("users").Filter(testCase.filter).First()
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}

	testCases = map[string]struct {
		expectedError error
		filter        map[string]any
	}{
		"nil params": {
			expectedError: ErrFilterParams,
			filter:        nil, // for nil params
		},
		"incorrect params": {
			expectedError: ErrRecordNotFound,
			filter:        map[string]any{"age": 0},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := fs.Memgodb().Collection("users").Filter(testCase.filter).First()
			require.ErrorIs(t, err, testCase.expectedError)
			assert.Nil(t, result)
		})
	}
}

func Test_Filter_All(t *testing.T) {
	fs := New()

	// insert a new records
	err := fs.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)

	testCases := map[string]struct {
		expectedError error
		filter        map[string]any
	}{
		"not nil params": {
			expectedError: nil,
			filter:        map[string]any{"age": 35.0}, // filter out records of age 35
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			results, err := fs.Memgodb().Collection("users").Filter(testCase.filter).All()
			require.NoError(t, err)
			assert.NotNil(t, results)
		})
	}

	testCases = map[string]struct {
		expectedError error
		filter        map[string]any
	}{
		"incorrect params": {
			expectedError: ErrRecordNotFound,
			filter:        map[string]any{"age": 0},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := fs.Memgodb().Collection("users").Filter(testCase.filter).All()
			require.ErrorIs(t, err, testCase.expectedError)
		})
	}
}

func Test_Delete_One(t *testing.T) {
	fs := New()

	// insert a new record
	err := fs.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)

	filters := map[string]map[string]any{
		"not nil params": {"age": 35.0}, // filter out record of age 35
	}

	for name, v := range filters {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("users").Delete(v).One()
			require.NoError(t, err)
		})
	}

	filters = map[string]map[string]any{
		"nil params": nil, // for nil params
	}

	for name, v := range filters {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("users").Delete(v).One()
			require.Error(t, err)
		})
	}
}

func Test_Delete_All(t *testing.T) {
	ch := Cache{}

	// insert a new record
	err := ch.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)

	filters := map[string]map[string]any{
		"not nil params": {"age": 35.0}, // filter out record of age 35
	}

	for name, v := range filters {
		t.Run(name, func(t *testing.T) {
			err := ch.Memgodb().Collection("users").Delete(v).All()
			require.NoError(t, err)
		})
	}

	filters = map[string]map[string]any{
		"nil params": nil, // for nil params
	}

	for name, v := range filters {
		t.Run(name, func(t *testing.T) {
			err := ch.Memgodb().Collection("users").Delete(v).All()
			require.NoError(t, err)
		})
	}
}

func Test_Update_One(t *testing.T) {
	fs := New()

	// insert a new record
	err := fs.Memgodb().Collection("user").Insert(MemgodbTestCases)
	require.NoError(t, err)

	testCases := map[string]struct {
		expectedError error
		filter        map[string]any
		update        map[string]any
	}{
		"correct filter params": {
			expectedError: nil,
			filter: map[string]any{
				"age": 35.0,
			},
			update: map[string]any{
				"age": 29,
			},
		},
		"nil filter params": {
			expectedError: ErrFilterParams,
			filter:        nil,
			update: map[string]any{
				"age": 28,
			},
		},
		"not found params": {
			expectedError: ErrRecordNotFound,
			filter: map[string]any{
				"age": 300.0,
			},
			update: map[string]any{
				"age": 28,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := fs.Memgodb().Collection("user").Update(testCase.filter, testCase.update).One()
			if err != nil {
				require.ErrorIs(t, err, testCase.expectedError)
			}
		})
	}
}

func Test_Persist(t *testing.T) {
	ch := Cache{}

	err := ch.Memgodb().Persist()
	require.NoError(t, err)
}

func Test_LoadDefault(t *testing.T) {
	ch := Cache{}

	err := ch.Memgodb().LoadDefault()
	require.NoError(t, err)
}
