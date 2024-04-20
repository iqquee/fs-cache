# fs-cache
A simple in-memory cache in golang


## Installation
```sh
go get github.com/iqquee/fs-cache@latest
```

## Usage

```go
package main

import (
	"fmt"

	fscache "github.com/iqquee/fs-cache"
)

func main() {
	fs := fscache.New()
    fs.Debug()

	if err := fs.Set("1", "user1", 5*time.Minute); err != nil {
		fmt.Println("set:", err)
	}

	if err := fs.Set("2", "user2"); err != nil {
		fmt.Println("error setting user2:", err)
	}

	if err := fs.Set("3", "user3"); err != nil {
		fmt.Println("error setting user3:", err)
	}

	if err := fs.Del("1"); err != nil {
		fmt.Println("error deleting key 1:", err)
	}

	result, err := fs.Get("2")
	if err != nil {
		fmt.Println("error getting key 2:", err)
	}

	fmt.Println("key 2 value:", result)

    if err := fs.Clear(); err != nil {
		fmt.Println("error clearing all datas:", err)
	}

	size := fs.Size()
	fmt.Println("total size: ", size)
}
```