package main

import (
	"./fs"
	"fmt"
	"os"
	"sort"
	"strings"
)

func PrintFolder(f Folder) {
	fmt.Print("files: ")
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
	for name := range f.Contents() {
		index = sort.Search(len(list), func(i int) bool {
			//behave more like ls
			if strings.ToLower(list[i]) == strings.ToLower(name) {
				return list[i] < name
			} else {
				return strings.ToLower(list[i]) > strings.ToLower(name)
			}
		})
		list = append(list, name)
		copy(list[index+1:], list[index:])
		list[index] = name
	}

	return list
}

func main() {
	var f Folder = fs.GetFolder(os.Args[1])

	update := make(chan bool)
	go f.Watch(update)

	for {
		select {
		case <-update:
			PrintFolder(f)
		}
	}

	done := make(chan bool)
	<-done
}
