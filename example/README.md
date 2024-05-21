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

# Memdis storage
Memdis gives you a Redis-like feature similarly as you would with a Redis database.
### Set()
Set() adds a new data into the in-memory storage
```go
fs := fscache.New()

// the third param is an optional param used to set the expiration time of the set data
if err := fs.Memdis().Set("key1", "user1", 5*time.Minute); err != nil {
	fmt.Println("error setting key1:", err)
}
```

### Get()
Get() retrieves a data from the in-memory storage
```go
fs := fscache.New()

result, err := fs.Memdis().Get("key1")
if err != nil {
	fmt.Println("error getting key 1:", err)
}

fmt.Println("key1:", result)
```

### SetMany()
SetMany() sets many data objects into memory for later access
```go
fs := fscache.New()

testCase := []map[string]fscache.MemdisData{
	{
		"key4": fscache.MemdisData{
			Value:    "value4",
			Duration: time.Now().Add(time.Minute),
		},
		"key5": fscache.MemdisData{
			Value: false,
		},
	},
}

setMany, err := fs.Memdis().SetMany(testCase)
if err != nil {
	fmt.Println("error setMany:", err)
}
fmt.Println("setMany:", setMany)
```

### GetMany()
GetMany() retrieves data with matching keys from the in-memory storage
```go
fs := fscache.New()

keys := []string{"key1", "key2"}

getMany := fs.Memdis().GetMany(keys)
fmt.Println("getMany:", getMany)
```

### OverWrite()
OverWrite() updates an already set value using it key
```go
fs := fscache.New()

if err := fs.Memdis().OverWrite("key1", "overwrite1", 1*time.Minute); err != nil {
	fmt.Println("error overwriting:", err)
}
```

### OverWriteWithKey()
OverWriteWithKey() updates an already set value and key using the previously set key
```go
fs := fscache.New()

if err := fs.Memdis().OverWriteWithKey("previousKey", "newKey", "newValue", 1*time.Minute); err != nil {
	fmt.Println("error overWriteWithKey:", err)
}
```

### Del()
Del() deletes a data from the in-memory storage
```go
fs := fscache.New()

if err := fs.Memdis().Del("key1"); err != nil {
	fmt.Println("error deleting key 1:", err)
}
```

### TypeOf()
TypeOf() returns the data type of a value
```go
fs := fscache.New()

typeOf, err := fs.Memdis().TypeOf("key1")
if err != nil {
	fmt.Println("error typeOf:", err)
}
fmt.Println("typeOf:", typeOf)
```

### Clear()
Clear() deletes all data from the in-memory storage
```go
fs := fscache.New()

if err := fs.Memdis().Clear(); err != nil {
	fmt.Println("error clearing all data:", err)
}
```

### Size()
Size() retrieves the total data objects in the in-memory storage
```go
fs := fscache.New()

size := fs.Memdis().Size()
fmt.Println("total size: ", size)
```

### Keys()
Keys() returns all the keys in the storage
```go
fs := fscache.New()

keys := fs.Memdis().Keys()
fmt.Println("keys: ", keys)
```

### Values()
Values() returns all the values in the storage
```go
fs := fscache.New()

values := fs.Memdis().Values()
fmt.Println("values: ", values)
```

### KeyValuePairs()
KeyValuePairs() returns an array of key value pairs of all the data in the storage
```go
fs := fscache.New()

keyValuePairs := fs.Memdis().KeyValuePairs()
fmt.Println("keyValuePairs: ", keyValuePairs)
```

# Memgodb storage
Memgodb gives you a MongoDB-like feature similarly as you would with a MondoDB database.

### Persist()
// Persist is used to write data to file. All data will be saved into a json file on the server.

This method will make sure all your your data's are saved into a json file. A cronJon runs ever minute and writes your data(s) into a json file to ensure data integrity

```go
fs := fscache.New()

if err := fs.Memgodb().Persist(); err != nil {
	fmt.Println(err)
}
```

### Insert()
Insert is used to insert a new record into the storage.
```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```
```go
fs := fscache.New()

// to insert single record
var user User
user.Name = "jane doe"
user.Age = 20

if err := fs.Memgodb().Collection(User{}).Insert(user); err != nil {
	fmt.Println(err)
}

// to insert multiple records
var users = []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
	{Name: "Jane doe",
		Age: 25},
	{Name: "Jane dice",
		Age: 20},
}

if err := fs.Memgodb().Collection("user").Insert(users); err != nil {
	fmt.Println(err)
}
```
### InsertFromJsonFile()
InsertFromJsonFile adds records into the storage from a JSON file.
```go
fs := fscache.New()

if err := fs.Memgodb().Collection("user").InsertFromJsonFile("path to JSON file"); err != nil {
	fmt.Println(err)
}
```

### Filter()
Filter is used to filter records from the storage. It has two methods which are First() and All().

- ### First()
First is a method available in Filter(), it returns the first matching record from the filter.

```go
type User struct {}
```
```go
fs := fscache.New()

// filter out record of age 35
filter := map[string]interface{}{
	"age": 35.0,
}

result, err := fs.Memgodb().Collection(User{}).Filter(filter).First()
if err != nil {
	fmt.Println(err)
}

fmt.Println(result)
```

- ### All()
All is a method available in Filter(), it returns the all matching records from the filter.

```go
type User struct {}
```
```go
fs := fscache.New()

// filter out record of age 35
filter := map[string]interface{}{
	"age": 35.0,
}

// to get all records with matching filter from the storage
matchingRecords, err := fs.Memgodb().Collection(User{}).Filter(filter).All()
if err != nil {
	fmt.Println(err)
}
fmt.Println(matchingRecords)

// to get all records from the collection from the storage
allRecords, err := fs.Memgodb().Collection(User{}).Filter(nil).All()
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
type User struct {}
```
```go
fs := fscache.New()

filter := map[string]interface{}{
	"age": 20.0,
}

if err := fs.Memgodb().Collection(User{}).Delete(filter).One(); err != nil {
	fmt.Println(err)
}
```
- ### Many()
Many adds many records into the storage at once
```go
type User struct {}
```

```go
fs := fscache.New()

filter := map[string]interface{}{
	"age": 20.0,
}

// to delete all records with matching filter from the storage
if err := fs.Memgodb().Collection(User{}).Delete(filter).All(); err != nil {
	fmt.Println(err)
}

// to delete all records in the collection from the storage
if err := fs.Memgodb().Collection(User{}).Delete(nil).All(); err != nil {
	fmt.Println(err)
}
```

### Update()
Update is used to update a existing record in the storage. It has a method which is One().

- ### One
One is a method available in Update(), it updates matching records from the filter, makes the necessry updated and returns an error if any.
```go
type User struct {}
```
```go
fs := fscache.New()

filter := map[string]interface{}{
	"age": 35.0,
}
update := map[string]interface{}{
	"age": 29,
}

if err := fs.Memgodb().Collection(User{}).Update(filter, update); err != nil {
	fmt.Println(err)
}
```

### LoadDefault
LoadDefault is used to load data from the json file saved on the server using Persist() if any.
```go
type User struct {}
```
```go
fs := fscache.New()

if err := fs.Memgodb().LoadDefault(); err != nil {
	fmt.Println(err)
}
```