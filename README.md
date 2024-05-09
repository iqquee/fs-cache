# fs-cache
fs-cache provides a quick way to store and retrieve frequently accessed data, significantly enhancing your application performance and reducing database queries / API calls.

## Features
- ### Key-value pair storage
- ### NoSql-like storage

## Installation
```sh
go get github.com/iqquee/fs-cache@v1.0.0 
```
if you wish to use the noSQL-like feature, then you should install like this
```sh
go get github.com/iqquee/fs-cache@latest
```

## Import
```sh
fscache "github.com/iqquee/fs-cache"
```

## Key-Value Pair Storage

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

## NoSql-like Storage

### Insert()
Insert adds a new record into the storage with collection name
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

if err := fs.NoSql().Collection(User{}).Insert(user); err != nil {
	fmt.Println(err)
}
```
### InsertMany()
InsertMany adds many records into the storage at once
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

if err := fs.NoSql().Collection("user").InsertMany(users); err != nil {
	fmt.Println(err)
}
```

### Filter()
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

// filter out records of age 35
filter := map[string]interface{}{
	"age": 35.0,
}

result, err := fs.NoSql().Collection(User{}).Filter(filter).All()
if err != nil {
	fmt.Println(err)
}

fmt.Println(result)
```

For an exhaustive documentation see the examples folder [https://github.com/iqquee/fs-cache/tree/main/example](https://github.com/iqquee/fs-cache/tree/main/example)

## Contributions
Anyone can contribute to this library ;). So, feel free to improve on and add new feaures. I await your pull requests.