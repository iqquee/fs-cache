package fscache

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type (
	// cacheData object
	cacheData struct {
		value    interface{}
		duration time.Time
	}

	// Cache object instance
	Cache struct {
		// debug enables debugging
		debug bool
		// Fscache is an [] that the datas are saved into
		Fscache []map[string]cacheData
	}

	// Operations lists all available operations on the fscache
	Operations interface {
		// Set() adds a new data into the in-memmory storage
		Set(key string, value interface{}, duration ...time.Duration) error
		// Get() retrieves a data from the in-memmory storage
		Get(key string) (interface{}, error)
		// Del() deletes a data from the in-memmory storage
		Del(key string) error
		// Clear() deletes all datas from the in-memmory storage
		Clear() error
		// Size() retrieves the total data objects in the in-memmory storage
		Size() int
		// Debug() enables debug to get certain logs
		Debug()
		// OverWrite updates an already set value using it key
		OverWrite(key string, value interface{}, duration ...time.Duration) error
		// OverWriteWithKey updates an already set value and key using the previously set key
		OverWriteWithKey(prevkey, newKey string, value interface{}, duration ...time.Duration) error
		// ExportJson exports all saves data objects as json
		ExportJson() []map[string]cacheData
		// ImportJson takes in an array of json objects and saves it into memory for later access
		ImportJson([]map[string]interface{}) error
		// Keys returns all the keys in the storage
		Keys() []string
		// Values returns all the values in the storage
		Values() []interface{}
		// TypeOf returns the data type of a value
		TypeOf(key string) (string, error)
		// SaveToFile saves the array of objects into a file
		SaveToFile(fileName string) error
	}
)

// New initializes an instance of the in-memory storage cache
func New() Operations {
	var fs []map[string]cacheData
	ch := Cache{
		Fscache: fs,
	}

	c := cron.New()

	// cron job set to run every 1 minute
	c.AddFunc("*/1 * * * *", func() {
		if ch.debug {
			fmt.Println("cron job running...")
		}
		for i := 0; i < len(ch.Fscache); i++ {
			for _, value := range ch.Fscache[i] {
				currenctTime := time.Now()
				if currenctTime.Before(value.duration) {
					if ch.debug {
						fmt.Printf("data object [%v] got expired ", ch.Fscache[i])
					}
					// take the data from off the array object
					ch.Fscache = append(ch.Fscache[:i], ch.Fscache[i+1:]...)
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
