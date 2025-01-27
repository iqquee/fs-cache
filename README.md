# fs-cache
fs-cache provides a quick way to store and retrieve frequently accessed data, significantly enhancing your application performance and reducing database queries / API calls.

## Features
- ### KeyStore storage
- ### DataStore storage

## Installation
```sh
go get github.com/iqquee/fs-cache@latest
```

## Import
```sh
fscache "github.com/iqquee/fs-cache"
```

## KeyStore storage
KeyStore gives you a Redis-like feature similarly as you would with a Redis database.

### Set()
Set() adds a new data into the in-memory storage
```go
fs := fscache.New()

// the third param is an optional param used to set the expiration time of the set data
if err := fs.KeyStore().Set("key1", "user1", 5*time.Minute); err != nil {
	fmt.Println("error setting key1:", err)
}
```

### Get()
Get() retrieves a data from the in-memory storage
```go
fs := fscache.New()

result, err := fs.KeyStore().Get("key1")
if err != nil {
	fmt.Println("error getting key 1:", err)
}

fmt.Println("key1:", result)
```

## DataStore storage
DataStore gives you an SQL/NoSQL-like feature.

### Create()
Create is used to insert a new record into the storage.
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

if err := fs.DataStore().Namespace(User{}).Create(user); err != nil {
	fmt.Println(err)
}
```

- ### First()
First is used when you expect a unique record and decodes it into the param object.
```go
fs := fscache.New()

// filter out record of age 35
filter := map[string]interface{}{
	"age": 35,
}

var response User
if err := fs.DataStore().Namespace(User{}).First(filter, &response); err != nil {
	fmt.Println(err)
}

fmt.Println(response)
```

- ### Find()
Find is used when you expect multiple records and decodes it into the param object.
```go
fs := fscache.New()

// filter out records with of age 35
filter := map[string]interface{}{
	"age": 35,
}

var response []User
if err := fs.DataStore().Namespace(User{}).Find(filter, &response); err != nil {
	fmt.Println(err)
}

fmt.Println(response)
```

- ### Sync()
You can use the Sync method to synchronize the records in the cache to your live sql database.
```go
fs := fscache.New()

// filter out records with of age 35
filter := map[string]interface{}{
	"age": 35,
}

ns := fs.DataStore().Namespace(User{})
db := gorm.Open(nil, &gorm.Config{})
ns.ConnectSQLDB(db).Sync(1 * time.Second)
```
<!-- 
For an exhaustive documentation see the examples folder [https://github.com/iqquee/fs-cache/tree/main/example](https://github.com/iqquee/fs-cache/tree/main/example) -->

## Contributions
Anyone can contribute to this library ;). So, feel free to improve on and add new features. I await your pull requests.