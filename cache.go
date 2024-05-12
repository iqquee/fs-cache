package fscache

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

var (
	debug bool
)

type (
	// MemdisData object
	MemdisData struct {
		Value    interface{}
		Duration time.Time
	}

	// Memdis object instance
	Memdis struct {
		logger zerolog.Logger
		// storage for key value pair storage
		storage []map[string]MemdisData
	}

	// Memgodb object instance
	Memgodb struct {
		logger zerolog.Logger
	}

	// Cache object
	Cache struct {
		MemdisInstance  Memdis
		MemgodbInstance Memgodb
	}

	// Operations lists all available operations on the fscache
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

	md := Memdis{
		logger:  logger,
		storage: memdicSorage,
	}

	Memgodb := Memgodb{
		logger: logger,
	}

	ch := Cache{
		MemdisInstance:  md,
		MemgodbInstance: Memgodb,
	}

	c := cron.New()

	// cron job set to run every 1 minute
	c.AddFunc("*/1 * * * *", func() {
		if debug {
			logger.Info().Msg("cron job running...")
		}

		if persistMemgodbData {
			if err := ch.MemgodbInstance.Persist(); err != nil {
				if debug {
					logger.Info().Msgf("persist error: %v", err)
				}
			}
		}

		for i := 0; i < len(ch.MemdisInstance.storage); i++ {
			for _, value := range ch.MemdisInstance.storage[i] {
				currenctTime := time.Now()
				if currenctTime.Before(value.Duration) {
					if debug {
						logger.Info().Msgf("data object [%v] got expired ", ch.MemdisInstance.storage[i])
					}
					// take the data from off the array object
					ch.MemdisInstance.storage = append(ch.MemdisInstance.storage[:i], ch.MemdisInstance.storage[i+1:]...)
					// decrement the array index by 1 since an object have been taken off the array
					i--
				}
			}
		}
	})

	c.Start()
	if debug {
		logger.Info().Msgf("cron job entries ::: %v", c.Entries())
	}

	op := Operations(&ch)
	return op
}

// Debug() enables debug to get certain logs
func (c *Cache) Debug() {
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
