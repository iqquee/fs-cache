## Installation
```sh
go get github.com/iqquee/fs-cache@v1.0.0 
```
## Import
```sh
fscache "github.com/iqquee/fs-cache"
```

### Debug()
Debug() enables debug to get certain logs
```go
fs := fscache.New()
// set if you want to get logs of activities
fs.Debug()
```

# Key-value pair storage

### Set()
Set() adds a new data into the in-memmory storage
```go
fs := fscache.New()

// the third param is an optional param used to set the expiration time of the set data
if err := fs.KeyValuePair().Set("key1", "user1", 5*time.Minute); err != nil {
	fmt.Println("error setting key1:", err)
}
```

### Get()
Get() retrieves a data from the in-memmory storage
```go
fs := fscache.New()

result, err := fs.KeyValuePair().Get("key1")
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

setMany, err := fs.KeyValuePair().SetMany(testCase)
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

getMany := fs.KeyValuePair().GetMany(keys)
fmt.Println("getMany:", getMany)
```

### OverWrite()
OverWrite() updates an already set value using it key
```go
fs := fscache.New()

if err := fs.KeyValuePair().OverWrite("key1", "overwrite1", 1*time.Minute); err != nil {
	fmt.Println("error overwriting:", err)
}
```

### OverWriteWithKey()
OverWriteWithKey() updates an already set value and key using the previously set key
```go
fs := fscache.New()

if err := fs.KeyValuePair().OverWriteWithKey("previousKey", "newKey", "newValue", 1*time.Minute); err != nil {
	fmt.Println("error overWriteWithKey:", err)
}
```

### Del()
Del() deletes a data from the in-memmory storage
```go
fs := fscache.New()

if err := fs.KeyValuePair().Del("key1"); err != nil {
	fmt.Println("error deleting key 1:", err)
}
```

### TypeOf()
TypeOf() returns the data type of a value
```go
fs := fscache.New()

typeOf, err := fs.KeyValuePair().TypeOf("key1")
if err != nil {
	fmt.Println("error typeOf:", err)
}
fmt.Println("typeOf:", typeOf)
```

### Clear()
Clear() deletes all datas from the in-memmory storage
```go
fs := fscache.New()

if err := fs.KeyValuePair().Clear(); err != nil {
	fmt.Println("error clearing all datas:", err)
}
```

### Size()
Size() retrieves the total data objects in the in-memmory storage
```go
fs := fscache.New()

size := fs.KeyValuePair().Size()
fmt.Println("total size: ", size)
```

### Keys()
Keys() returns all the keys in the storage
```go
fs := fscache.New()

keys := fs.KeyValuePair().Keys()
fmt.Println("keys: ", keys)
```

### Values()
Values() returns all the values in the storage
```go
fs := fscache.New()

values := fs.KeyValuePair().Values()
fmt.Println("values: ", values)
```

### KeyValuePairs()
KeyValuePairs() returns an array of key value pairs of all the datas in the storage
```go
fs := fscache.New()

keyValuePairs := fs.KeyValuePair().KeyValuePairs()
fmt.Println("keyValuePairs: ", keyValuePairs)
```

# NoSql-like storage

### Insert()
Insert is used to insert a new record into the storage. It has two methods which are One() and Many().

- ### One
One adds a new record into the storage with collection name
```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```
```go
fs := fscache.New()

var user User
user.Name = "jane doe" 
user.Age = 20

res, err := fs.NoSql().Collection(User{}).Insert(user).One
if err != nil {
	fmt.Println(err)
}

fmt.Println(res)
```
- ### Many()
Many adds many records into the storage at once
```go
fs := fscache.New()

var users = []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
	{Name: "Jane doe",
		Age: 25},
	{Name: "Jane dice",
		Age: 20},
}

if err := fs.NoSql().Collection("user").Insert(nil).Many(users); err != nil {
	fmt.Println(err)
}
```

### Filter()
Filter is used to filter records from the storage. It has two methods which are First() and All().

- ### First()
First is a method available in Filter(), it returns the first matching record from the filter.
```go
fs := fscache.New()

// filter out record of age 35
filter := map[string]interface{}{
	"age": 35.0,
}

result, err := fs.NoSql().Collection(User{}).Filter(filter).First()
if err != nil {
	fmt.Println(err)
}

fmt.Println(result)
```

- ### All()
All is a method available in Filter(), it returns the all matching records from the filter.
```go
fs := fscache.New()

// filter out record of age 35
filter := map[string]interface{}{
	"age": 35.0,
}

// to get all records with matching filter from the storage
matchingRecords, err := fs.NoSql().Collection(User{}).Filter(filter).All()
if err != nil {
	fmt.Println(err)
}
fmt.Println(matchingRecords)

// to get all records from the collection from the storage
allRecords, err := fs.NoSql().Collection(User{}).Filter(nil).All()
if err != nil {
	fmt.Println(err)
}

fmt.Println(allRecords)
```

### Delete()
Delete is used to delete a new record from the storage. It has two methods which are One() and Many().

- ### One
One is a method available in Delete(), it deletes a record and returns an error if any.
```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```
```go
fs := fscache.New()

filter := map[string]interface{}{
	"age": 20.0,
}

if err := fs.NoSql().Collection("user").Delete(filter).One(); err != nil {
	fmt.Println(err)
}
```
- ### Many()
Many adds many records into the storage at once
```go
fs := fscache.New()

filter := map[string]interface{}{
	"age": 20.0,
}

// to delete all records with matching filter from the storage
if err := fs.NoSql().Collection("user").Delete(filter).All(); err != nil {
	fmt.Println(err)
}

// to delete all records in the collection from the storage
if err := fs.NoSql().Collection("user").Delete(nil).All(); err != nil {
	fmt.Println(err)
}
```