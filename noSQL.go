package fscache

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	// Collection object
	Collection struct {
		logger         zerolog.Logger
		collectionName string
		storage        []interface{}
	}

	Result struct{}
)

// Collection defines the collection(table) name to perform an operations on
func (ns *NoSQL) Collection(col interface{}) *Collection {
	t := reflect.TypeOf(col)

	// run validation
	if reflect.ValueOf(col).IsZero() && col == nil {
		ns.logger.Error().Msg("Collection cannot be empty...")
		panic("Collection cannot be empty...")
	}

	if t.Kind() != reflect.Struct && t.Kind() != reflect.String {
		ns.logger.Error().Msg("Collection must either be a [string] or an [object]")
		panic("Collection must either be a [string] or an [object]")
	}

	var colName string
	if t.Kind() == reflect.Struct {
		colName = strings.ToLower(t.Name())
	} else {
		colName = strings.ToLower(col.(string))
	}

	// check if the last index ends with the letter s
	// if not, append (s) to it e.g user = users
	if len(colName) > 0 && string(colName[len(colName)-1]) != "s" {
		colName = fmt.Sprintf("%ss", colName)
	}

	return &Collection{
		logger:         ns.logger,
		collectionName: colName,
		storage:        ns.storage,
	}
}

// Insert adds a new record into the storage with collection name
func (c *Collection) Insert(obj interface{}) (interface{}, error) {
	t := reflect.TypeOf(obj)

	if t.Kind() != reflect.Struct && t.Kind() != reflect.Map {
		c.logger.Error().Msg("Functions params must either be a [map] or a [struct]")
		panic("Functions params must either be a [map] or an [struct]")
	}

	objMap := make(map[string]interface{})
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		c.logger.Err(err).Msg("JSON marshal error")
		return nil, err
	}

	if err := json.Unmarshal(jsonObj, &objMap); err != nil {
		c.logger.Err(err).Msg("JSON unmarshal error")
		return nil, err
	}

	// add additional data to the object
	objMap["colName"] = c.collectionName
	objMap["id"] = uuid.New()
	objMap["createdAt"] = time.Now()
	objMap["deletedAt"] = nil

	c.storage = append(c.storage, objMap)
	return objMap, nil
}

// InsertMany adds many record into the storage at once
func (c *Collection) InsertMany(arr interface{}) error {
	t := reflect.TypeOf(arr)

	if t.Kind() != reflect.Slice {
		c.logger.Error().Msg("Functions params must must be an [slice]")
		return errors.New("functions params must be an [slice]")
		// panic("Functions params must be an [slice]")
	}

	var arrObjs []map[string]interface{}
	arrObj, err := json.Marshal(arr)
	if err != nil {
		c.logger.Err(err).Msg("JSON marshal error")
		return err
	}

	if err := json.Unmarshal(arrObj, &arrObjs); err != nil {
		c.logger.Err(err).Msg("JSON unmarshal error")
		return err
	}

	for _, obj := range arrObjs {
		objMap := make(map[string]interface{})
		jsonObj, err := json.Marshal(obj)
		if err != nil {
			c.logger.Err(err).Msg("JSON marshal error")
			return err
		}

		if err := json.Unmarshal(jsonObj, &objMap); err != nil {
			c.logger.Err(err).Msg("JSON unmarshal error")
			return err
		}

		// add additional data to the object
		objMap["colName"] = c.collectionName
		objMap["id"] = uuid.New()
		objMap["createdAt"] = time.Now()
		objMap["deletedAt"] = nil

		c.storage = append(c.storage, objMap)
	}

	fmt.Println(c.storage)
	return nil
}
