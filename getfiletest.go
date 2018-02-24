package main

import (
	"./fm"
	"./fs"
	"fmt"
	"os"
	"sync"
)

var printLock sync.Mutex

func PrintFile(f fm.File) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Println(f.Name())
	fmt.Println(f.Size())
	fmt.Println(f.ModTime())
	fmt.Println("Directory: ", f.IsDir())
	fmt.Print("\n")

}

func main() {
	f, err := fs.Stat(os.Args[1])
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	PrintFile(f)
}
