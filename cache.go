package fscache

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var debug bool

type (
	// MemdisData object
	MemdisData struct {
		Value    any
		Duration time.Time
	}

	// Memdis object instance
	Memdis struct {
		mu     *sync.RWMutex
		logger zerolog.Logger
		// storage for key value pair storage
		storage []map[string]MemdisData
	}

	// Memgodb object instance
	Memgodb struct {
		mu     *sync.RWMutex
		logger zerolog.Logger
	}

	// Cache object
	Cache struct {
		logger          zerolog.Logger
		MemdisInstance  Memdis
		MemgodbInstance Memgodb
	}

	// Operations lists all available operations on the fs-cache
	Operations interface {
		// Debug() enables debug to get certain logs
		Debug()

		// Memdis gives you a Redis-like feature similarly as you would with a Redis database
		Memdis() *Memdis
		// Memgodb gives you a MongoDB-like feature similarly as you would with a MondoDB database
		Memgodb() *Memgodb
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	var memdicSorage []map[string]MemdisData
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	mu := sync.RWMutex{}
	md := Memdis{
		mu:      &mu,
		logger:  logger,
		storage: memdicSorage,
	}

	Memgodb := Memgodb{
		mu:     &mu,
		logger: logger,
	}

	ch := Cache{
		logger:          logger,
		MemdisInstance:  md,
		MemgodbInstance: Memgodb,
	}

	// start go routine
	go ch.runner()

	op := Operations(&ch)
	return op
}

// Debug() enables debug to get certain logs
func (*Cache) Debug() {
	debug = true
}

// KeyValue returns methods for key-value pair storage
func (c *Cache) Memdis() *Memdis {
	return &c.MemdisInstance
}

// Memgodb returns methods for Memgodb-like storage
func (c *Cache) Memgodb() *Memgodb {
	return &Memgodb{
		logger: c.MemgodbInstance.logger,
	}
}

// runner runs every 30 seconds to persists the Memgodb records and delete expired records from the Memdis storage.
func (ch *Cache) runner() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if debug {
			ch.logger.Info().Msg("cron job running...")
		}

		if persistMemgodbData {
			if err := ch.MemgodbInstance.Persist(); err != nil {
				if debug {
					ch.logger.Info().Msgf("persist error: %v", err)
				}
			}
		}

		for i := 0; i < len(ch.MemdisInstance.storage); i++ {
			for _, value := range ch.MemdisInstance.storage[i] {
				currentTime := time.Now()
				if currentTime.Before(value.Duration) {
					ch.Memdis().mu.Lock()
					if debug {
						ch.logger.Info().Msgf("data object [%v] got expired ", ch.MemdisInstance.storage[i])
					}
					// take the data from off the array object
					ch.MemdisInstance.storage = append(ch.MemdisInstance.storage[:i], ch.MemdisInstance.storage[i+1:]...)
					// decrement the array index by 1 since an object have been taken off the array
					i--
					ch.Memdis().mu.Unlock()
				}
			}
		}
	}
}
