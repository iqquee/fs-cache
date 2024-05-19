package fscache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var (
	// MemgodbStorage storage instance
	MemgodbStorage []any
	// persistMemgodbData to enable persistence of Memgodb data
	persistMemgodbData bool
)

var (
	// ErrRecordNotFound record not found
	ErrRecordNotFound = errors.New("record not found")
	// ErrFilterParams filter params cannot be nil
	ErrFilterParams = errors.New("filter params cannot be nil")
)

type (
	// Collection object
	Collection struct {
		logger         zerolog.Logger
		collectionName string
	}

	// Insert object implements One() and Many() to insert new records
	Insert struct {
		obj        any
		collection Collection
	}

	// Filter object implements One() and All()
	Filter struct {
		objMaps    []map[string]any
		filter     map[string]any
		collection Collection
	}

	// Delete object implements One() and All()
	Delete struct {
		objMaps    []map[string]any
		filter     map[string]any
		collection Collection
	}

	// Persist objects implemented Persist() used to persist inserted records
	Persist struct {
		Error error
	}

	// Update object implements One() and All()
	Update struct {
		objMaps    []map[string]any
		filter     map[string]any
		update     map[string]any
		collection Collection
	}
)

// Collection defines the collection(table) name to perform an operation on it
func (mg *Memgodb) Collection(col any) *Collection {
	t := reflect.TypeOf(col)

	// run validation
	if reflect.ValueOf(col).IsZero() && col == nil {
		mg.logger.Error().Msg("Collection cannot be empty...")
		panic("Collection cannot be empty...")
	}

	if t.Kind() != reflect.Struct && t.Kind() != reflect.String {
		mg.logger.Error().Msg("Collection must either be a [string] or an [object]")
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
		logger:         mg.logger,
		collectionName: colName,
	}
}

// Insert is used to insert a new record into the storage. It has two methods which are One() and Many().
func (c *Collection) Insert(obj any) *Insert {
	return &Insert{
		obj:        obj,
		collection: *c,
	}
}

// One is a method available in Insert(). It adds a new record into the storage with collection name
func (i *Insert) One() (any, error) {
	if i.obj == nil {
		return nil, errors.New("One() params cannot be nil")
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
	objMap["updatedAt"] = nil

	MemgodbStorage = append(MemgodbStorage, objMap)
	return objMap, nil
}

// Many is a method available in Insert(). It adds many records into the storage at once
func (i *Insert) Many(arr any) ([]any, error) {
	if i.obj != nil {
		return nil, errors.New("Many() params must be nil to insert Many")
	}

	t := reflect.TypeOf(arr)

	if t.Kind() != reflect.Slice {
		return nil, errors.New("function param must be a [slice]")
	}

	arrObjs, err := i.collection.decodeMany(arr)
	if err != nil {
		return nil, err
	}

	var savedData []any
	for _, obj := range arrObjs {
		saved, err := i.collection.Insert(obj).One()
		if err != nil {
			return nil, err
		}

		savedData = append(savedData, saved)
	}

	return savedData, nil
}

// FromJsonFile is a method available in Insert(). It adds records into the storage from a JSON file
func (i *Insert) FromJsonFile(fileLocation string) error {
	if i.obj != nil {
		return errors.New("FromFile() params must be nil to insert from file")
	}

	f, err := os.Open(fileLocation)
	if err != nil {
		return err
	}
	defer f.Close()

	fileByte, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var obj any
	if err := json.Unmarshal(fileByte, &obj); err != nil {
		return errors.New("invalid json file")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Slice {
		objMap, err := i.collection.decodeMany(obj)
		if err != nil {
			return nil
		}

		_, err = i.collection.Insert(nil).Many(objMap)
		if err != nil {
			return nil
		}
	} else if t.Kind() == reflect.Map {
		objMap, err := i.collection.decode(obj)
		if err != nil {
			return nil
		}

		_, err = i.collection.Insert(objMap).One()
		if err != nil {
			return nil
		}
	} else {
		return errors.New("file must contain either an array of [objects ::: slice] or [object ::: map]")
	}

	return nil
}

// Filter is used to filter records from the storage. It has two methods which are First() and All().
func (c *Collection) Filter(filter map[string]any) *Filter {
	var objMaps []map[string]any
	var err error

	if filter != nil {
		objMaps, err = c.decodeMany(MemgodbStorage)
		if err != nil {
			return nil
		}
	}

	return &Filter{
		objMaps:    objMaps,
		filter:     filter,
		collection: *c,
	}
}

// First is a method available in Filter(), it returns the first matching record from the filter.
func (f *Filter) First() (map[string]any, error) {
	if f.objMaps == nil {
		return nil, ErrFilterParams
	}

	notFound := true
	var foundObj map[string]any
	counter := 0
	for _, item := range f.objMaps {
		for key, val := range f.filter {
			if item["colName"] == f.collection.collectionName {
				if v, ok := item[key]; ok && val == v {
					if counter < 1 {
						notFound = false
						foundObj = item
						counter++
					}
					break
				}
			}
		}
	}

	if notFound {
		return nil, ErrRecordNotFound
	}

	return foundObj, nil
}

// All is a method available in Filter(), it returns all the matching records from the filter.
func (f *Filter) All() ([]map[string]any, error) {
	if f.objMaps == nil {
		var objMaps []map[string]any
		arrObj, err := json.Marshal(MemgodbStorage)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(arrObj, &objMaps); err != nil {
			return nil, err
		}

		return objMaps, nil
	}

	notFound := true
	var foundObj []map[string]any
	for _, item := range f.objMaps {
		for key, val := range f.filter {
			if item["colName"] == f.collection.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					foundObj = append(foundObj, item)
				}
			}
		}
	}

	if notFound {
		return nil, ErrRecordNotFound
	}

	return foundObj, nil
}

// Delete is used to delete a new record from the storage. It has two methods which are One() and Many().
func (c *Collection) Delete(filter map[string]any) *Delete {
	var objMaps []map[string]any
	var err error

	if filter != nil {
		objMaps, err = c.decodeMany(MemgodbStorage)
		if err != nil {
			return nil
		}
	}

	return &Delete{
		objMaps:    objMaps,
		filter:     filter,
		collection: *c,
	}
}

// One is a method available in Delete(), it deletes a record and returns an error if any.
func (d *Delete) One() error {
	if d.objMaps == nil {
		return ErrFilterParams
	}

	notFound := true
	for index, item := range d.objMaps {
		for key, val := range d.filter {
			if item["colName"] == d.collection.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					if index < (len(MemgodbStorage) - 1) {
						MemgodbStorage = append(MemgodbStorage[:index], MemgodbStorage[index+1:]...)
						index--
						break
					} else {
						MemgodbStorage = MemgodbStorage[:index]
						break
					}
				}
			}
		}
	}

	if notFound {
		return ErrRecordNotFound
	}

	return nil
}

// All is a method available in Delete(), it deletes matching records from the filter and returns an error if any.
func (d *Delete) All() error {
	if d.objMaps == nil {
		MemgodbStorage = MemgodbStorage[:0]
		return nil
	}

	notFound := true
	for index, item := range d.objMaps {
		for key, val := range d.filter {
			if item["colName"] == d.collection.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					if index < (len(MemgodbStorage) - 1) {
						MemgodbStorage = append(MemgodbStorage[:index], MemgodbStorage[index+1:]...)
						index--
					} else {
						MemgodbStorage = MemgodbStorage[:index]
					}
				}
			}
		}
	}

	if notFound {
		return ErrRecordNotFound
	}

	return nil
}

// Update is used to update an existing record in the storage. It has a method which is One().
func (c *Collection) Update(filter, obj map[string]any) *Update {
	var objMaps []map[string]any
	var err error

	if filter != nil {
		objMaps, err = c.decodeMany(MemgodbStorage)
		if err != nil {
			return nil
		}
	}

	return &Update{
		objMaps:    objMaps,
		filter:     filter,
		update:     obj,
		collection: *c,
	}
}

// One is a method available in Update(), it updates matching records from the filter, makes the necessary updates and returns an error if any.
func (u *Update) One() error {
	if u.objMaps == nil {
		return ErrFilterParams
	}

	notFound := true
	counter := 0
	for index, item := range u.objMaps {
		for key, val := range u.filter {
			if item["colName"] == u.collection.collectionName {
				if v, ok := item[key]; ok && val == v {
					notFound = false
					if counter < 1 {
						for _, updateValue := range u.update {
							item[key] = updateValue
							counter++
							break
						}
						item["updatedAt"] = time.Now()
					}
					MemgodbStorage[index] = item
				}
			}
		}
	}

	if notFound {
		return ErrRecordNotFound
	}

	return nil
}

// LoadDefault is used to load data from the JSON file saved on the server using Persist() if any.
func (mg *Memgodb) LoadDefault() error {
	f, err := os.Open("./memgodbstorage.json")
	if err != nil {
		return errors.New("error finding file")
	}
	defer f.Close()

	fileByte, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var obj any
	if err := json.Unmarshal(fileByte, &obj); err != nil {
		return errors.New("invalid json file")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Slice {
		var objMap []any
		jsonByte, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonByte, &objMap); err != nil {
			return err
		}

		MemgodbStorage = append(MemgodbStorage, objMap...)
	} else if t.Kind() == reflect.Map {
		var objMap any
		jsonByte, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonByte, &objMap); err != nil {
			return err
		}

		MemgodbStorage = append(MemgodbStorage, objMap)
	}

	return nil
}

// Persist is used to write data to file. All data will be saved into a JSON file on the server.

// This method will make sure all your data are saved into a JSON file. A cron job runs ever minute and writes your data into a JSON file to ensure data integrity
func (mg *Memgodb) Persist() error {
	if MemgodbStorage == nil {
		return nil
	}

	persistMemgodbData = true
	jsonByte, err := json.Marshal(MemgodbStorage)
	if err != nil {
		return err
	}

	file, err := os.Create("./memgodbstorage.json")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonByte)
	if err != nil {
		return err
	}

	return nil
}

// decode decodes an any into a map[string]any
func (*Collection) decode(obj any) (map[string]any, error) {
	objMap := make(map[string]any)
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonObj, &objMap); err != nil {
		return nil, err
	}

	return objMap, nil
}

// decodeMany decodes an any into an []map[string]any
func (*Collection) decodeMany(arr any) ([]map[string]any, error) {
	var arrObjs []map[string]any
	arrObj, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(arrObj, &arrObjs); err != nil {
		return nil, err
	}

	return arrObjs, nil
}
