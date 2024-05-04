package fscache

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type (
	// cacheData object
	CacheData struct {
		Value    interface{}
		Duration time.Time
	}

	// KeyPair object instance
	KeyPair struct {
		// Storage for key value pair storage
		Storage []map[string]CacheData
	}

	// NoSQL object instance
	NoSQL struct {
		// Storage for NoSQL-like storage
		Storage []map[string]interface{}
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
	var noSql []map[string]interface{}

	kp := KeyPair{
		Storage: keyValuePair,
	}

	noSQL := NoSQL{
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
			fmt.Println("cron job running...")
		}
		for i := 0; i < len(ch.KeyPair.Storage); i++ {
			for _, value := range ch.KeyPair.Storage[i] {
				currenctTime := time.Now()
				if currenctTime.Before(value.Duration) {
					if ch.debug {
						fmt.Printf("data object [%v] got expired ", ch.KeyPair.Storage[i])
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
		fmt.Printf("cron job entries ::: %v", c.Entries())
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
// client.NoSQL.Collection(struct{}).InsertMany(struct{})
// filter
// client.NoSQL.Collection(struct{}).Filter([]map[string]interface{})
// update
// client.NoSQL.Collection(struct{}).Filter([]map[string]interface{}).Update(struct{})
// delete
// client.NoSQL.Collection(struct{}).Delete([]map[string]interface{})
// filter and return all
// client.NoSQL.Collection(struct{}).Filter([]map[string]interface{}).All()
// filter and return all paginated
// client.NoSQL.Collection(struct{}).Filter([]map[string]interface{}).All().Paginate(pd ...pageDetails)
