package main

import (
	"./fm"
	"./fs"
	"fmt"
	"os"
	"sync"
	"time"
)

var printLock sync.Mutex

func PrintFolder(f fm.Folder) {
	printLock.Lock()
	defer printLock.Unlock()
	for _, file := range fm.GetSortedFileList(f) {
		fmt.Print(file, " ")
	}
	fmt.Print("\n")
}

func main() {
	/*
		for i := 0; i < 1000; i++ {
			go watchFolder(os.Args[1])
		}
	*/

	for _, arg := range os.Args[2:] {
		go watchFolder(arg)
	}

	watchFolder(os.Args[1])

	fmt.Println("no folders")
	time.Sleep(5 * time.Second) //observe cleanup
	fmt.Println("we are done here")
}

func watchFolder(path string) {
	fmt.Println(path)
	f, err := fs.GetFolder(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	update := f.Watch()

Loop:
	for {
		select {
		case <-update:
			if f.Path() == "" && f.Contents() == nil {
				fmt.Println("folder object gone")
				break Loop
			}
			PrintFolder(f)

			//weird mechanism for testing this
			for file := range f.Contents() {
				if file == "close" {
					f.Close()
					break Loop
				}
			}
		}
	}
}
