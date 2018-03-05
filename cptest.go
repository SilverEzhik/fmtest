package main

import (
	"./fs"
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) != 2 {
		fmt.Println("need two arguments")
		return
	}

	err := fs.CopyFile(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
	}
}
