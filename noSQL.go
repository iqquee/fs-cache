package fscache

import (
	"encoding/json"
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
func (c *Collection) InsertMany(objs []interface{}) error {
	for _, obj := range objs {
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

	return nil
}
