# fs-cache
fs-cache provides a quick way to store and retrieve frequently accessed data, significantly enhancing your application performance and reducing database queries / API calls.

## Features
- ### Memdis storage
- ### Memgodb storage

## Installation
```sh
go get github.com/iqquee/fs-cache@latest
```

## Import
```sh
fscache "github.com/iqquee/fs-cache"
```

## Memdis storage
Memdis gives you a Redis-like feature similarly as you would with a Redis database.

### Set()
Set() adds a new data into the in-memmory storage
```go
fs := fscache.New()

// the third param is an optional param used to set the expiration time of the set data
if err := fs.Memdis().Set("key1", "user1", 5*time.Minute); err != nil {
	fmt.Println("error setting key1:", err)
}
```

### Get()
Get() retrieves a data from the in-memmory storage
```go
fs := fscache.New()

result, err := fs.Memdis().Get("key1")
if err != nil {
	fmt.Println("error getting key 1:", err)
}

fmt.Println("key1:", result)
```

## Memgodb storage
Memgodb gives you a MongoDB-like feature similarly as you would with a MondoDB database.

### Persist()
Persist is used to write data to file. All datas will be saved into a json file.

This method will make sure all your your data's are saved. A cronJon runs ever minute and writes your data(s) into a json file to ensure data integrity
```go
fs := fscache.New()

if err := fs.Memgodb().Persist(); err != nil {
	fmt.Println(err)
}
```

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

res, err := fs.Memgodb().Collection(User{}).Insert(user).One()
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

if err := fs.Memgodb().Collection("user").Insert(nil).Many(users); err != nil {
	fmt.Println(err)
}
```
- ### FromJsonFile()
FromJsonFile is a method available in Insert(). It adds record(s) into the storage from a json file
```go
fs := fscache.New()

if err := fs.Memgodb().Collection("user").Insert(nil).FromJsonFile("path to JSON file"); err != nil {
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

result, err := fs.Memgodb().Collection(User{}).Filter(filter).First()
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


For an exhaustive documentation see the examples folder [https://github.com/iqquee/fs-cache/tree/main/example](https://github.com/iqquee/fs-cache/tree/main/example)

## Contributions
Anyone can contribute to this library ;). So, feel free to improve on and add new feaures. I await your pull requests.