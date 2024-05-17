package fscache

import (
	"context"
	"os"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { // launch "cron job"
		for {
			// wait for context cancellation
			select {
			case <-ctx.Done():
				return

			// wait 1 minute before continuing
			case <-time.After(1 * time.Minute):
			}

			// cron job set to run every 1 minute
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
		}
	}()

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
