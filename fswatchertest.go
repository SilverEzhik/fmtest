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

	queue := false
	for {
		select {
		case <-update:
			if f.Path() == "" && f.Contents() == nil {
				fmt.Println("folder object gone")
				break
			}

			// It doesn't do much good to update on every single event, especially since
			// removing many files, for example, will send out individual events.
			// So we limit to only updating every 100 ms.
			// At the same time, this should be done on the UI layer, because different apps
			// might want different rate limits and such.
			go func() {
				if queue == false {
					queue = true
					go PrintFolder(f)
					time.Sleep(100 * time.Millisecond)
					queue = false
				}
			}()

		}
	}
}
