package fscache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	// Collection object
	Collection struct {
		logger         zerolog.Logger
		CollectionName string
		Storage        []interface{}
	}

	Result struct{}
)

// Collection defines the collection(table) name to perform an operations on
func (ns *NoSQL) Collection(col interface{}) *Collection {
	t := reflect.TypeOf(col)
	// validate t.Name()
	if len(t.Name()) == 0 {
		ns.logger.Error().Msg("Collection cannot be empty...")
	}

	if t.Kind() != reflect.Struct || t.Kind() != reflect.String {
		ns.logger.Error().Msg("Collection must either be a [string] or an [object]")
	}

	var colName string
	// check if the last index ends with the letter s
	// if not, append (s) to it e.g user = users
	if len(t.Name()) > 0 && string(t.Name()[len(t.Name())-1]) != "s" {
		colName = fmt.Sprintf("%ss", t.Name())
	} else {
		colName = t.Name()
	}

	return &Collection{
		logger:         ns.logger,
		CollectionName: colName,
		Storage:        ns.Storage,
	}
}

// Insert inserts a new record into the storage with collection name
func (c *Collection) Insert(obj interface{}) (interface{}, error) {
	var objMap map[string]interface{}

	v := reflect.TypeOf(obj)
	if v.Kind() != reflect.Map {
		jsonObj, err := json.Marshal(obj)
		if err != nil {
			c.logger.Err(err).Msg("JSON marshal error")
			return nil, err
		}

		if err := json.Unmarshal(jsonObj, &objMap); err != nil {
			c.logger.Err(err).Msg("JSON unmarshal error")
			return nil, err
		}
	}

	// add additional data to the object
	objMap["id"] = uuid.New()
	objMap["createdAt"] = time.Now()

	c.Storage = append(c.Storage, objMap)
	return obj, nil
}
