package fscache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	DBTypePostgres DBType = "postgres"
	DBTypeMysql    DBType = "mysql"
)

type (
	DBType    string
	Namespace struct {
		dataStore *DataStore
		namespace string
	}

	Find  struct{}
	First struct{}

	ConnectSQLDB struct {
		DBType    DBType
		namespace *Namespace
		DB        *gorm.DB
	}

	SQLDBConfig struct {
		DBType       DBType
		IsConnString bool
		ConnString   string
		Address      string
		Port         string
		User         string
		Database     string
		Password     string
	}

	ConnectMongoDB struct{}
)

// Namespace creates or retrieves a namespace within the DataStore.
// If a schema is provided, it will be associated with the namespace.
// If no schema is provided, the namespace will be initialized with a nil schema.
// The function returns a Namespace struct containing the logger, data, indexes, schemas, and mutex from the DataStore.
//
// Parameters:
//   - name: The name of the namespace to create or retrieve.
//   - schemas: Optional variadic parameter to provide a schema for the namespace.
//
// Returns:
//   - Namespace: A struct containing the logger, data, indexes, schemas, and mutex from the DataStore.
func (ds *DataStore) Namespace(name any, schema ...Schema) Namespace {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	t := reflect.TypeOf(name)
	if reflect.ValueOf(name).IsZero() && name == nil {
		ds.logger.Error().Msg("Namespace cannot be empty.")
		panic("Error ::: Namespace cannot be empty.")
	}

	if t.Kind() != reflect.Struct && t.Kind() != reflect.String {
		ds.logger.Error().Msg("Namespace must either be a [string] or a [struct]")
		panic("Error ::: Namespace must either be a [string] or a [struct]")
	}

	var nameSpace string
	if t.Kind() == reflect.Struct {
		nameSpace = strings.ToLower(t.Name())
	} else {
		nameSpace = strings.ToLower(name.(string))
	}

	if len(nameSpace) > 0 && string(nameSpace[len(nameSpace)-1]) != "s" {
		nameSpace = fmt.Sprintf("%ss", nameSpace)
	}

	// If no schema is passed, initialize it as nil or empty
	var schemas Schema
	if len(schema) > 0 {
		schemas = schema[0] // Use the first schema if passed
	} else {
		schemas = nil // No schema provided
	}

	ds.schemas[nameSpace] = schemas
	ds.indexes[nameSpace] = make(map[string]map[any][]int)

	return Namespace{
		dataStore: ds,
		namespace: nameSpace,
	}
}

// toSnakeCase converts a given CamelCase string to snake_case.
// It inserts an underscore before each uppercase letter (except the first one)
// and converts all characters to lowercase.
//
// Parameters:
//
//	s - the CamelCase string to be converted.
//
// Returns:
//
//	A snake_case representation of the input string.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_', r)
		} else {
			result = append(result, r)
		}
	}

	return strings.ToLower(string(result))
}

// Create adds a new entry to the namespace's data store. It first locks the data store
// to ensure thread safety. If a schema is defined for the namespace, it enforces the schema
// by checking the types of the provided values. If the types do not match the schema, it returns an error.
// After validation, it appends the entry to the data store and updates the indexes for quick lookups.
//
// Parameters:
//
//	v - A map containing the key-value pairs to be added to the data store.
//
// Returns:
//
//	error - An error if the schema validation fails, otherwise nil.
func (ns *Namespace) Create(v map[string]any) error {
	ns.dataStore.mu.Lock()
	defer ns.dataStore.mu.Unlock()

	normalized := make(map[string]any)
	for key, val := range v {
		normalized[toSnakeCase(key)] = val
	}

	// Schema enforcement
	if schema, ok := ns.dataStore.schemas[ns.namespace]; ok {
		for key, val := range normalized {
			if expectedType, exists := schema[key]; exists {
				if reflect.TypeOf(val).String() != expectedType {
					ns.dataStore.logger.Error().Msgf("Error ::: invalid type for field %s: expected %s, got %s", key, expectedType, reflect.TypeOf(val).String())
					return fmt.Errorf("invalid type for field %s: expected %s, got %s", key, expectedType, reflect.TypeOf(val).String())
				}
			}
		}
	}

	// Add a field of isSynced to each record inserted
	normalized["is_synced"] = false
	ns.dataStore.data[ns.namespace] = append(ns.dataStore.data[ns.namespace], normalized)

	// Update indexes
	for key, value := range normalized {
		if _, exists := ns.dataStore.indexes[ns.namespace][key]; !exists {
			ns.dataStore.indexes[ns.namespace][key] = make(map[any][]int)
		}
		ns.dataStore.indexes[ns.namespace][key][value] = append(ns.dataStore.indexes[ns.namespace][key][value], len(ns.dataStore.data[ns.namespace])-1)
	}

	return nil
}

// Query retrieves documents from the namespace's data store that match the provided filters.
// It returns a slice of maps, where each map represents a document, and an error if any occurs.
//
// Parameters:
//
//	filters - A map where the key is the field name and the value is the value to filter by.
//
// Returns:
//
//	A slice of maps, where each map represents a document that matches the filters.
//	An error if any occurs during the query process.
func (ns *Namespace) Query(filters map[string]any) ([]map[string]any, error) {
	var result []map[string]any

	if len(filters) == 0 {
		return ns.dataStore.data[ns.namespace], nil
	}

	docIndexes := make(map[int]bool)

	for key, value := range filters {
		if idx, exists := ns.dataStore.indexes[ns.namespace][toSnakeCase(key)]; exists {
			if docIdxs, exists := idx[value]; exists {
				for _, idx := range docIdxs {
					docIndexes[idx] = true
				}
			}
		}
	}

	for idx := range docIndexes {
		result = append(result, ns.dataStore.data[ns.namespace][idx])
	}

	return result, nil
}

// Find searches for records in the namespace that match the given filters.
// If one result is found, an error is returned suggesting to use First() for one result.
// The results are decoded into the provided variable.
//
// Parameters:
//   - filters: A map of filter criteria to apply to the query.
//   - v: A variable to store the decoded results.
//
// Returns:
//   - error: An error if the query fails, if only one result is found, or if decoding fails.
func (ns *Namespace) Find(filters map[string]any, v any) error {
	result, err := ns.Query(filters)
	if err != nil {
		return err
	}

	if err := ns.decodeMany(result, &v); err != nil {
		return err
	}

	return nil
}

// First retrieves the first result matching the provided filters and decodes it into the provided variable.
// If more than one result is found, an error is returned suggesting to use Find() for multiple results.
// Only use First() if the expected result is a single object.
// Parameters:
//
//	filters - a map containing the filters to apply to the query.
//	v - a variable to decode the result into.
//
// Returns:
//
//	error - an error if the query fails, more than one result is found, or decoding fails.
func (ns *Namespace) First(filters map[string]any, v any) error {
	result, err := ns.Query(filters)
	if err != nil {
		return err
	}

	if len(result) > 1 {
		return fmt.Errorf("find() expects one result, but got %d. Use Find() instead for multiple data", len(result))
	}

	if err := ns.decodeOne(result[0], &v); err != nil {
		return err
	}

	return nil
}

// decodeMany takes a slice of maps with string keys and values of any type,
// marshals it into JSON, and then unmarshals it into the provided variable v.
// It returns an error if either the marshaling or unmarshalling process fails.
//
// Params:
// - params: A slice of maps where each map has string keys and values of any type.
// - v: A pointer to a variable where the unmarshaled JSON will be stored.
//
// Returns:
// - error: An error if marshaling or unmarshalling fails, otherwise nil.
func (ns *Namespace) decodeMany(params []map[string]any, v any) error {
	jsobByte, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsobByte, &v); err != nil {
		return err
	}

	return nil
}

// decodeOne decodes a map of parameters into a given value.
//
// Parameters:
//   - params: A map containing the parameters to be decoded.
//   - v: A pointer to the value where the decoded parameters will be stored.
//
// Returns:
//   - An error if the encoding or decoding process fails, otherwise nil.
func (ns *Namespace) decodeOne(params map[string]any, v any) error {
	jsobByte, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsobByte, &v); err != nil {
		return err
	}

	return nil
}

// Update modifies documents in the data store that match the given filters with the provided new data.
// It acquires a lock on the data store to ensure thread safety, queries for matching documents,
// updates each matching document with the new data, and rebuilds indexes if necessary.
//
// Parameters:
//   - filters: A map of key-value pairs used to filter documents that need to be updated.
//   - newData: A map of key-value pairs representing the new data to be applied to the matching documents.
//
// Returns:
//   - error: An error if the query fails or any other issue occurs during the update process.
func (ns *Namespace) Update(filters map[string]any, newData map[string]any) error {
	ns.dataStore.mu.Lock()
	defer ns.dataStore.mu.Unlock()

	matchingDocs, err := ns.Query(filters)
	if err != nil {
		return err
	}

	for _, doc := range matchingDocs {
		for key, value := range newData {
			doc[toSnakeCase(key)] = value
		}
	}

	// Rebuild indexes if necessary
	ns.rebuildIndexes()

	return nil
}

// cs.namespace.dataStore.indexes[namespace]["isSynced"][false] = append(cs.namespace.dataStore.indexes[namespace]["isSynced"][true], index)
// Delete removes documents from the namespace's data store that match the given filters.
// It first queries the data store to find matching documents, then removes them from the slice,
// and finally rebuilds the indexes.
//
// Parameters:
//
//	filters - a map of field names to values that documents must match to be deleted.
//
// Returns:
//
//	error - an error if the query fails, otherwise nil.
func (ns *Namespace) Delete(filters map[string]any) error {
	ns.dataStore.mu.Lock()
	defer ns.dataStore.mu.Unlock()

	// Perform query first to find matching documents
	matchingDocs, err := ns.Query(filters)
	if err != nil {
		return err
	}

	// Remove matching documents from the slice
	for _, doc := range matchingDocs {
		for i, storedDoc := range ns.dataStore.data[ns.namespace] {
			if reflect.DeepEqual(storedDoc, doc) {
				// Delete the document from the slice
				ns.dataStore.data[ns.namespace] = append(ns.dataStore.data[ns.namespace][:i], ns.dataStore.data[ns.namespace][i+1:]...)
				break
			}
		}
	}

	// Rebuild indexes after deletion
	ns.rebuildIndexes()

	return nil
}

// rebuildIndexes rebuilds the indexes for the namespace.
// It resets the current index and iterates over all documents
// in the namespace to recreate the index based on the document
// keys and values. Each value is mapped to a list of document
// indices where it appears.
func (ns *Namespace) rebuildIndexes() {
	// Reset the namespace index
	ns.dataStore.indexes[ns.namespace] = make(map[string]map[any][]int)

	// Iterate over all documents in the namespace
	for i, doc := range ns.dataStore.data[ns.namespace] {
		for key, value := range doc {
			if _, exists := ns.dataStore.indexes[ns.namespace][key]; !exists {
				ns.dataStore.indexes[ns.namespace][key] = make(map[any][]int)
			}

			ns.dataStore.indexes[ns.namespace][key][value] = append(ns.dataStore.indexes[ns.namespace][key][value], i)
		}
	}
}

// ListNamespaces returns a list of all namespace names present in the DataStore.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (ds *DataStore) ListNamespaces() []string {
	var namespaces []string
	for namespace := range ds.schemas {
		namespaces = append(namespaces, namespace)
	}

	return namespaces
}

// ConnectSQLDB establishes a connection to the SQL database using the provided gorm.DB instance.
// It returns a pointer to a ConnectSQLDB struct that contains the database connection and the namespace.
//
// Parameters:
//
//	db - A pointer to a gorm.DB instance representing the database connection.
//
// Returns:
//
//	A pointer to a ConnectSQLDB struct containing the database connection and namespace.
func (ns *Namespace) ConnectSQLDB(db *gorm.DB) *ConnectSQLDB {
	// if db == nil {
	// 	return nil, fmt.Errorf("invalid database connection")
	// }

	return &ConnectSQLDB{DB: db, namespace: ns}
}

// Sync periodically synchronizes the in-memory data store with the SQL database.
// It runs in a separate goroutine and uses a ticker to trigger the synchronization
// process at the specified interval.
//
// Parameters:
//   - interval: The duration between each synchronization attempt.
//
// The synchronization process involves the following steps:
//  1. Lock the data store to ensure thread safety.
//  2. Iterate over each namespace and its associated records.
//  3. For each document, check if it is already synced by looking at the "is_synced" field.
//  4. If the document is not synced, create a copy of the document excluding the "is_synced" field.
//  5. Attempt to sync the document to the SQL database.
//  6. If the sync is successful, mark the document as synced and log the success.
//  7. If an error occurs during the sync, log the error and continue with the next document.
//  8. Unlock the data store after processing all documents.
//
// Note: The synchronization process continues indefinitely until the program terminates.
func (cs *ConnectSQLDB) Sync(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			cs.namespace.dataStore.mu.Lock()
			for namespace, records := range cs.namespace.dataStore.data {
				for index, doc := range records {
					if isSynced, ok := doc["is_synced"].(bool); ok && isSynced {
						continue
					}

					docCopy := make(map[string]any)
					for key, value := range doc {
						if key != "is_synced" {
							docCopy[key] = value
						}
					}

					if err := cs.DB.Table(namespace).Create(docCopy).Error; err != nil {
						cs.namespace.dataStore.logger.Err(err).Msgf(
							"Error syncing document at index %d in namespace %s: %v",
							index, namespace, err,
						)
						// Log the error and continue with the next document
						continue
					}

					doc["is_synced"] = true
					cs.namespace.dataStore.data[namespace][index] = doc

					cs.namespace.dataStore.logger.Info().Msgf(
						"Synced document at index %d in namespace %s",
						index, namespace,
					)
				}
			}
			cs.namespace.dataStore.mu.Unlock()
		}
	}()
}

func (ds *DataStore) ConnectMongoDB() error {
	return nil
}

func (ds *ConnectMongoDB) Sync() error {
	return nil
}
