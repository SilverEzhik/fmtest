package main

import (
	"./fm"
	"./fs"
	"fmt"
	"os"
	"time"
)

func PrintFolder(f fm.Folder) {
	for _, file := range fm.GetSortedFileList(f) {
		fmt.Print(file, " ")
	}
	fmt.Print("\n")
}

func main() {
	watchFolder()
}

func watchFolder() {
	for _, path := range os.Args[1:] {
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
		time.Sleep(5 * time.Second) //observe cleanup
		fmt.Println("we are done here")
	}
}
