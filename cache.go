package fscache

import (
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type (
	// KeyStoreData object
	KeyStoreData struct {
		Value    any
		Duration time.Time
	}

	// KeyStore object instance
	KeyStore struct {
		mu     *sync.RWMutex
		logger zerolog.Logger
		// storage for key value pair storage
		storage []map[string]KeyStoreData
	}

	// DataStore represents the in-memory store for documents (key-value pairs)
	DataStore struct {
		logger  zerolog.Logger
		data    map[string][]map[string]interface{}         // Map to store a slice of documents per namespace
		indexes map[string]map[string]map[interface{}][]int // Indexes for fast querying
		schemas map[string]Schema                           // Schema for validation
		mu      *sync.RWMutex                               // Mutex for thread safety
	}

	// Schema represents the structure of a document with type validation
	Schema map[string]string

	// Cache object
	Cache struct {
		logger            zerolog.Logger
		KeyStoreInstance  KeyStore
		DataStoreInstance DataStore
	}

	// Operations lists all available operations on the fs-cache
	Operations interface {
		// Debug() enables debug to get certain logs
		Debug(io.Writer)

		// KeyStore gives you a Redis-like feature similarly as you would with a Redis database
		KeyStore() *KeyStore
		// DataStore gives you a MongoDB-like feature similarly as you would with a MondoDB database
		DataStore() *DataStore
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	logger := zerolog.New(io.Discard)
	mu := &sync.RWMutex{}

	ks := KeyStore{
		mu:      mu,
		logger:  logger,
		storage: make([]map[string]KeyStoreData, 0),
	}

	ds := DataStore{
		logger:  logger,
		mu:      mu,
		data:    make(map[string][]map[string]interface{}),
		indexes: make(map[string]map[string]map[interface{}][]int),
		schemas: make(map[string]Schema),
	}

	ch := Cache{
		logger:            logger,
		KeyStoreInstance:  ks,
		DataStoreInstance: ds,
	}

	// start go routine
	go ch.runner()

	op := Operations(&ch)
	return op
}

// Debug() enables debug to get certain logs
func (c *Cache) Debug(w io.Writer) {
	logger := zerolog.New(w).With().Timestamp().Logger()
	c.logger = logger
	c.KeyStoreInstance.logger = logger
	c.DataStoreInstance.logger = logger
}

// KeyStore returns methods for key-value pair storage
func (c *Cache) KeyStore() *KeyStore {
	return &c.KeyStoreInstance
}

// DataStore returns methods for a NoSQL or SQL-[like] storage
func (c *Cache) DataStore() *DataStore {
	return &c.DataStoreInstance
}

// runner is a method of the Cache struct that periodically performs maintenance tasks.
// It runs a cron job every 30 seconds to:
// 1. Log the execution of the cron job.
// 2. Persist data if the persistDataStoreData flag is set.
// 3. Lock the KeyStore, check for expired data objects, and remove them from the storage.
//
// The method uses a ticker to trigger the cron job at regular intervals and ensures
// that the ticker is stopped when the method exits.
func (ch *Cache) runner() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ch.logger.Info().Msg("cron job running...")

		// Persist data if necessary
		if persistDataStoreData {
			if err := ch.DataStoreInstance.Persist(); err != nil {
				ch.logger.Info().Msgf("persist error: %v", err)
			}
		}

		ch.KeyStore().mu.Lock()
		defer ch.KeyStore().mu.Unlock()

		var toRemove []int

		currentTime := time.Now()

		for i := 0; i < len(ch.KeyStoreInstance.storage); i++ {
			for _, value := range ch.KeyStoreInstance.storage[i] {
				if currentTime.Before(value.Duration) {
					toRemove = append(toRemove, i)
					ch.logger.Info().Msgf("data object [%v] got expired", ch.KeyStoreInstance.storage[i])
					break
				}
			}
		}

		for _, index := range toRemove {
			ch.KeyStoreInstance.storage = append(ch.KeyStoreInstance.storage[:index], ch.KeyStoreInstance.storage[index+1:]...)
		}
	}
}
