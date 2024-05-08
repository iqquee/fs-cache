# fs-cache
fs-cache provides a quick way to store and retrieve frequently accessed data, significantly enhancing your application performance and reducing database queries / API calls. It supports data caching with an optional params for expiration

## Features
- ### Key-value pair storage
<!-- - ### NoSql-like storage -->

## Installation
```sh
go get github.com/iqquee/fs-cache@v1.0.0 
```
## Import
```sh
fscache "github.com/iqquee/fs-cache"
```

## Usage

### Debug()
Debug() enables debug to get certain logs
```go
fs := fscache.New()
// set if you want to get logs of activities
fs.Debug()
```

### Set()
Set() adds a new data into the in-memmory storage
```go
fs := fscache.New()

// the third param is an optional param used to set the expiration time of the set data
if err := fs.Set("key1", "user1", 5*time.Minute); err != nil {
	fmt.Println("error setting key1:", err)
}
```

### Get()
Get() retrieves a data from the in-memmory storage
```go
fs := fscache.New()

result, err := fs.Get("key1")
if err != nil {
	fmt.Println("error getting key 1:", err)
}

fmt.Println("key1:", result)
```

### SetMany()
SetMany() sets many data objects into memory for later access
```go
fs := fscache.New()

testCase := []map[string]fscache.CacheData{
	{
		"key4": fscache.CacheData{
			Value:    "value4",
			Duration: time.Now().Add(time.Minute),
		},
		"key5": fscache.CacheData{
			Value: false,
		},
	},
}

setMany, err := fs.SetMany(testCase)
if err != nil {
	fmt.Println("error setMany:", err)
}
fmt.Println("setMany:", setMany)
```

### GetMany()
GetMany() retrieves datas with matching keys from the in-memmory storage
```go
fs := fscache.New()

keys := []string{"key1", "key2"}

getMany := fs.GetMany(keys)
fmt.Println("getMany:", getMany)
```

### OverWrite()
OverWrite() updates an already set value using it key
```go
fs := fscache.New()

if err := fs.OverWrite("key1", "overwrite1", 1*time.Minute); err != nil {
	fmt.Println("error overwriting:", err)
}
```

### OverWriteWithKey()
OverWriteWithKey() updates an already set value and key using the previously set key
```go
fs := fscache.New()

if err := fs.OverWriteWithKey("previousKey", "newKey", "newValue", 1*time.Minute); err != nil {
	fmt.Println("error overWriteWithKey:", err)
}
```

### Del()
Del() deletes a data from the in-memmory storage
```go
fs := fscache.New()

if err := fs.Del("key1"); err != nil {
	fmt.Println("error deleting key 1:", err)
}
```

### TypeOf()
TypeOf() returns the data type of a value
```go
fs := fscache.New()

typeOf, err := fs.TypeOf("key1")
if err != nil {
	fmt.Println("error typeOf:", err)
}
fmt.Println("typeOf:", typeOf)
```

### Clear()
Clear() deletes all datas from the in-memmory storage
```go
fs := fscache.New()

if err := fs.Clear(); err != nil {
	fmt.Println("error clearing all datas:", err)
}
```

### Size()
Size() retrieves the total data objects in the in-memmory storage
```go
fs := fscache.New()

size := fs.Size()
fmt.Println("total size: ", size)
```

### Keys()
Keys() returns all the keys in the storage
```go
fs := fscache.New()

keys := fs.Keys()
fmt.Println("keys: ", keys)
```

### Values()
Values() returns all the values in the storage
```go
fs := fscache.New()

values := fs.Values()
fmt.Println("values: ", values)
```

### KeyValuePairs()
KeyValuePairs() returns an array of key value pairs of all the datas in the storage
```go
fs := fscache.New()

keyValuePairs := fs.KeyValuePairs()
fmt.Println("keyValuePairs: ", keyValuePairs)
```