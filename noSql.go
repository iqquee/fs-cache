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

// noSql storage instance
var noSqlStorage []interface{}

type (
	// Collection object
	Collection struct {
		logger         zerolog.Logger
		collectionName string
	}

	// Insert object implementes One() and Many() to insert new records
	Insert struct {
		obj        interface{}
		collection Collection
	}

	// Filter object implementes One() and All()
	Filter struct {
		objMaps        []map[string]interface{}
		filter         map[string]interface{}
		collectionName string
	}

	// Delete object implementes One() and All()
	Delete struct {
		objMaps        []map[string]interface{}
		filter         map[string]interface{}
		collectionName string
	}

	// Persist objects implemented Persist() used to persist inserted records
	Persist struct{}
)

// Collection defines the collection(table) name to perform an operations on
func (ns *NoSQL) Collection(col interface{}) *Collection {
	t := reflect.TypeOf(col)

	// run validation
	if reflect.ValueOf(col).IsZero() && col == nil {
		if debug {
			ns.logger.Error().Msg("Collection cannot be empty...")
		}
		panic("Collection cannot be empty...")
	}

	if t.Kind() != reflect.Struct && t.Kind() != reflect.String {
		if debug {
			ns.logger.Error().Msg("Collection must either be a [string] or an [object]")
		}
		panic("Collection must either be a [string] or an [object]")
	}

	var colName string
	if t.Kind() == reflect.Struct {
		colName = strings.ToLower(t.Name())
	} else {
		colName = strings.ToLower(col.(string))
	}

	if len(colName) > 0 && string(colName[len(colName)-1]) != "s" {
		colName = fmt.Sprintf("%ss", colName)
	}

	return &Collection{
		logger:         ns.logger,
		collectionName: colName,
	}
}

// Insert is used to insert a new record into the storage. It has two methods which are One() and Many().
func (c *Collection) Insert(obj interface{}) *Insert {
	return &Insert{
		obj:        obj,
		collection: *c,
	}
}

// One is a method available in Insert(). It adds a new record into the storage with collection name
func (i *Insert) One() (interface{}, error) {
	if i.obj == nil {
		return nil, errors.New("insert() params cannot be nil")
	}

	t := reflect.TypeOf(i.obj)

	if t.Kind() != reflect.Struct && t.Kind() != reflect.Map {
		return nil, errors.New("insert() param must either be a [map] or a [struct]")
	}

	objMap, err := i.collection.decode(i.obj)
	if err != nil {
		return nil, err
	}

	objMap["colName"] = i.collection.collectionName
	objMap["id"] = uuid.New()
	objMap["createdAt"] = time.Now()
	objMap["deletedAt"] = nil

	noSqlStorage = append(noSqlStorage, objMap)
	return objMap, nil
}

// Many is a method available in Insert(). It adds many records into the storage at once
func (i *Insert) Many(arr interface{}) error {
	if i.obj != nil {
		return errors.New("insert() params must be nil to insert Many")
	}

	t := reflect.TypeOf(arr)

	if t.Kind() != reflect.Slice {
		return errors.New("function param must be a [slice]")
	}

	arrObjs, err := i.collection.decodeMany(arr)
	if err != nil {
		return err
	}

	for _, obj := range arrObjs {
		if _, err := i.collection.Insert(obj).One(); err != nil {
			return err
		}
	}

	return nil
}

// Filter is used to filter records from the storage. It has two methods which are First() and All().
func (c *Collection) Filter(filter map[string]interface{}) *Filter {
	var objMaps []map[string]interface{}
	var err error

	if filter != nil {
		objMaps, err = c.decodeMany(noSqlStorage)
		if err != nil {
			return nil
		}
	}

	return &Filter{
		objMaps:        objMaps,
		filter:         filter,
		collectionName: c.collectionName,
	}
}

// First is a method available in Filter(), it returns the first matching record from the filter.
func (f *Filter) First() (map[string]interface{}, error) {
	if f.objMaps == nil {
		return nil, errors.New("filter params cannot be nil")
	}

	notFound := true
	var foundObj map[string]interface{}
	for _, item := range f.objMaps {
		for key, val := range f.filter {
			if item["colName"] == f.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					foundObj = item
					break
				}
			}
		}
	}

	if notFound {
		return nil, errors.New("record not found")
	}

	return foundObj, nil
}

// All is a method available in Filter(), it returns all the matching records from the filter.
func (f *Filter) All() ([]map[string]interface{}, error) {
	if f.objMaps == nil {
		var objMaps []map[string]interface{}
		arrObj, err := json.Marshal(noSqlStorage)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(arrObj, &objMaps); err != nil {
			return nil, err
		}

		return objMaps, nil
	}

	notFound := true
	var foundObj []map[string]interface{}
	for _, item := range f.objMaps {
		for key, val := range f.filter {
			if item["colName"] == f.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					foundObj = append(foundObj, item)
				}
			}
		}
	}

	if notFound {
		return nil, errors.New("record not found")
	}

	return foundObj, nil
}

// Delete is used to delete a new record from the storage. It has two methods which are One() and Many().
func (c *Collection) Delete(filter map[string]interface{}) *Delete {
	var objMaps []map[string]interface{}
	var err error

	if filter != nil {
		objMaps, err = c.decodeMany(noSqlStorage)
		if err != nil {
			return nil
		}
	}

	return &Delete{
		objMaps:        objMaps,
		filter:         filter,
		collectionName: c.collectionName,
	}
}

// One is a method available in Delete(), it deletes a record and returns an error if any.
func (d *Delete) One() error {
	if d.objMaps == nil {
		return errors.New("filter params cannot be nil")
	}

	notFound := true
	for index, item := range d.objMaps {
		for key, val := range d.filter {
			if item["colName"] == d.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					if index < (len(noSqlStorage) - 1) {
						noSqlStorage = append(noSqlStorage[:index], noSqlStorage[index+1:]...)
						index--
						break
					} else {
						noSqlStorage = noSqlStorage[:index]
						break
					}
				}
			}
		}
	}

	if notFound {
		return errors.New("record not found")
	}

	return nil
}

// All is a method available in Delete(), it deletes matching records from the filter and returns an error if any.
func (d *Delete) All() error {
	if d.objMaps == nil {
		noSqlStorage = noSqlStorage[:0]
		return nil
	}

	notFound := true
	for index, item := range d.objMaps {
		for key, val := range d.filter {
			if item["colName"] == d.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					fmt.Println("Delected: ", item)
					if index < (len(noSqlStorage) - 1) {
						noSqlStorage = append(noSqlStorage[:index], noSqlStorage[index+1:]...)
						index--
					} else {
						noSqlStorage = noSqlStorage[:index]
					}
				}
			}
		}
	}

	if notFound {
		return errors.New("record not found")
	}

	return nil
}

// decode decodes an interface{} into a map[string]interface{}
func (c *Collection) decode(obj interface{}) (map[string]interface{}, error) {
	objMap := make(map[string]interface{})
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonObj, &objMap); err != nil {
		return nil, err
	}

	return objMap, nil
}

// decodeMany decodes an interface{} into an []map[string]interface{}
func (c *Collection) decodeMany(arr interface{}) ([]map[string]interface{}, error) {
	var arrObjs []map[string]interface{}
	arrObj, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(arrObj, &arrObjs); err != nil {
		return nil, err
	}

	return arrObjs, nil
}
