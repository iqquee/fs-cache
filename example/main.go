package main

import (
	"fmt"

	fscache "github.com/iqquee/fs-cache"
)

func main() {
	fs := fscache.New()
	if err := fs.Set("1", "user1", 10*5); err != nil {
		fmt.Println("set:", err)
	}

	if err := fs.Set("2", "user2"); err != nil {
		fmt.Println("set:", err)
	}

	if err := fs.Set("3", "user3"); err != nil {
		fmt.Println("set:", err)
	}

	if err := fs.Del("1"); err != nil {
		fmt.Println("del:", err)
	}

	// if err := fs.Clear(); err != nil {
	// 	fmt.Println("clear:", err)
	// }

	result, err := fs.Get("2")
	if err != nil {
		fmt.Println("get:", err)
	}

	fmt.Println(result)

	size := fs.Size()
	fmt.Println("Size: ", size)
}
