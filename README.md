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
	// set if you want to get logs of activities
    fs.Debug()

	if err := fs.Set("1", "user1", 5*time.Minute); err != nil {
		fmt.Println("set:", err)
	}

	if err := fs.Del("1"); err != nil {
		fmt.Println("error deleting key 1:", err)
	}

	result, err := fs.Get("1")
	if err != nil {
		fmt.Println("error getting key 1:", err)
	}

	fmt.Println("key 1 value:", result)

    if err := fs.Clear(); err != nil {
		fmt.Println("error clearing all datas:", err)
	}

	size := fs.Size()
	fmt.Println("total size: ", size)
}
```