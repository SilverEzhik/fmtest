package main

import (
	"./fs"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func PrintFolder(f Folder) {
	for _, file := range GetSortedFileList(f) {
		fmt.Print(file, " ")
	}
	fmt.Print("\n")
}

//Map is unsorted, which is not as fun.
//Would need other sorting mechanisms, probably just keep this away from the actual folder object.
//not very efficient doing this every time?
func GetSortedFileList(f Folder) []string {
	list := make([]string, 0, len(f.Contents()))

	index := 0

	//https://gist.github.com/zhum/57cb45d8bbea86d87490
	for _, file := range f.Contents() {
		index = sort.Search(len(list), func(i int) bool {
			//behave more like ls
			if strings.ToLower(list[i]) == strings.ToLower(file.Name()) {
				return list[i] < file.Name()
			} else {
				return strings.ToLower(list[i]) > strings.ToLower(file.Name())
			}
		})
		list = append(list, file.Name())
		copy(list[index+1:], list[index:])
		list[index] = file.Name()
	}

	return list
}

func main() {
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
