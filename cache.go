package fscache

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type (
	// cacheData object
	CacheData struct {
		Value    interface{}
		Duration time.Time
	}

	// KeyPair object instance
	KeyPair struct {
		logger zerolog.Logger
		// Storage for key value pair storage
		Storage []map[string]CacheData
	}

	// NoSQL object instance
	NoSQL struct {
		logger zerolog.Logger
		// Storage for NoSQL-like storage
		Storage []interface{}
	}

	// Cache object
	Cache struct {
		// debug enables debugging
		debug   bool
		KeyPair KeyPair
		NoSQL   NoSQL
	}

	// Operations lists all available operations on the fscache
	Operations interface {
		// Debug() enables debug to get certain logs
		Debug()

		KeyValuePair() *KeyPair
		NoSql() *NoSQL
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	var keyValuePair []map[string]CacheData
	var noSql []interface{}
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	kp := KeyPair{
		logger:  logger,
		Storage: keyValuePair,
	}

	noSQL := NoSQL{
		logger:  logger,
		Storage: noSql,
	}

	ch := Cache{
		KeyPair: kp,
		NoSQL:   noSQL,
	}

	c := cron.New()

	// cron job set to run every 1 minute
	c.AddFunc("*/1 * * * *", func() {
		if ch.debug {
			logger.Info().Msg("cron job running...")
		}
		for i := 0; i < len(ch.KeyPair.Storage); i++ {
			for _, value := range ch.KeyPair.Storage[i] {
				currenctTime := time.Now()
				if currenctTime.Before(value.Duration) {
					if ch.debug {
						logger.Info().Msgf("data object [%v] got expired ", ch.KeyPair.Storage[i])
					}
					// take the data from off the array object
					ch.KeyPair.Storage = append(ch.KeyPair.Storage[:i], ch.KeyPair.Storage[i+1:]...)
					// decrement the array index by 1 since an object have been taken off the array
					i--
				}
			}
		}
	})

	c.Start()
	if ch.debug {
		logger.Info().Msgf("cron job entries ::: %v", c.Entries())
	}

	op := Operations(&ch)
	return op
}

// KeyValue returns KeyPair object
func (c *Cache) KeyValuePair() *KeyPair {
	return &c.KeyPair
}

// NoSql returns NoSql object
func (c *Cache) NoSql() *NoSQL {
	return &c.NoSQL
}

// create
// client.NoSQL.Collection(struct{}).Insert(struct{})

// create many
// client.NoSQL.Collection(struct{}).InsertMany(struct{})

// filter
// client.NoSQL.Collection(struct{}).Find(map[string]interface{})

// filter many
// client.NoSQL.Collection(struct{}).FindMany([]map[string]interface{})

// update
// client.NoSQL.Collection(struct{}).Find([]map[string]interface{}).Update(struct{})

// delete one
// client.NoSQL.Collection(struct{}).DeleteOne(map[string]interface{})

// delete many
// client.NoSQL.Collection(struct{}).DeleteMany([]map[string]interface{})

// filter and return all
// client.NoSQL.Collection(struct{}).Find([]map[string]interface{}).All()

// filter and return all paginated
// client.NoSQL.Collection(struct{}).Find([]map[string]interface{}).All().Paginate(pd ...pageDetails)

// count
// client.NoSQL.Collection(struct{}).Count()

// delete all datas in a collection
// client.NoSQL.Collection(struct{}).Remove()

// deletea all collections data
// client.NoSQL.DropAll()
