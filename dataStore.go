package fscache

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Namespace struct {
		dataStore *DataStore
		namespace string
	}
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
		ds.logger.Error().Msg("Namespace must either be a [string] or an [object]")
		panic("Error ::: Namespace must either be a [string] or an [object]")
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
	ds.indexes[nameSpace] = make(map[string]map[interface{}][]int)

	return Namespace{
		dataStore: ds,
		namespace: nameSpace,
	}
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
func (ns *Namespace) Create(v map[string]interface{}) error {
	ns.dataStore.mu.Lock()
	defer ns.dataStore.mu.Unlock()

	// Schema enforcement
	if schema, ok := ns.dataStore.schemas[ns.namespace]; ok {
		for key, val := range v {
			if expectedType, exists := schema[key]; exists {
				if reflect.TypeOf(val).String() != expectedType {
					return fmt.Errorf("invalid type for field %s: expected %s, got %s", key, expectedType, reflect.TypeOf(val).String())
				}
			}
		}
	}

	ns.dataStore.data[ns.namespace] = append(ns.dataStore.data[ns.namespace], v)

	for key, value := range v {
		if _, exists := ns.dataStore.indexes[ns.namespace][key]; !exists {
			ns.dataStore.indexes[ns.namespace][key] = make(map[interface{}][]int)
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
func (ns *Namespace) Query(filters map[string]interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	docIndexes := make(map[int]bool)

	for key, value := range filters {
		if idx, exists := ns.dataStore.indexes[ns.namespace][key]; exists {
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
func (ns *Namespace) Update(filters map[string]interface{}, newData map[string]interface{}) error {
	ns.dataStore.mu.Lock()
	defer ns.dataStore.mu.Unlock()

	matchingDocs, err := ns.Query(filters)
	if err != nil {
		return err
	}

	for _, doc := range matchingDocs {
		for key, value := range newData {
			doc[key] = value
		}
	}

	// Rebuild indexes if necessary
	ns.rebuildIndexes()

	return nil
}

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
func (ns *Namespace) Delete(filters map[string]interface{}) error {
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
	ns.dataStore.indexes[ns.namespace] = make(map[string]map[interface{}][]int)

	// Iterate over all documents in the namespace
	for i, doc := range ns.dataStore.data[ns.namespace] {
		for key, value := range doc {
			if _, exists := ns.dataStore.indexes[ns.namespace][key]; !exists {
				ns.dataStore.indexes[ns.namespace][key] = make(map[interface{}][]int)
			}

			ns.dataStore.indexes[ns.namespace][key][value] = append(ns.dataStore.indexes[ns.namespace][key][value], i)
		}
	}
}

// ListNamespaces returns a list of all namespace names present in the DataStore.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (ds *DataStore) ListNamespaces() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var namespaces []string
	for namespace := range ds.schemas {
		namespaces = append(namespaces, namespace)
	}

	return namespaces
}
