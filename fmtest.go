package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Mutex it
type Folder struct {
	path     string
	contents map[string]os.FileInfo
}

// Get initial folder contents from file system
// Might be nice to have a mechanism that would just go ahead and refresh if
// Many update/remove events are queued.
func (f *Folder) Refresh() {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f.contents[file.Name()] = file
		//fmt.Println(file.Name())
	}
}

// Takes an absolute path to file and stats it
// Also functions as an add function
func (f *Folder) Update(path string) {
	file, err := os.Stat(path)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	f.contents[filepath.Base(path)] = file
	//fmt.Printf("up: %s\n", f.contents[filepath.Base(path)].Name())
}

func (f *Folder) Remove(path string) {
	delete(f.contents, filepath.Base(path))
}

func (f *Folder) Print() {
	fmt.Print("files: ")
	for _, file := range f.GetSortedFileList() {
		fmt.Print(file, " ")
	}
	fmt.Print("\n")
}

//Map is unsorted, which is not as fun.
//Would need other sorting mechanisms, probably just keep this away from the actual folder object.
//not very efficient doing this every time?
func (f *Folder) GetSortedFileList() []string {
	list := make([]string, 0, len(f.contents))

	index := 0

	//https://gist.github.com/zhum/57cb45d8bbea86d87490
	for name := range f.contents {
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

	pwd := Folder{path: os.Args[1], contents: make(map[string]os.FileInfo)}
	pwd.Refresh()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				//create or modify - update map
				case fsnotify.Create:
					fallthrough
				case fsnotify.Rename:
					fallthrough
				case fsnotify.Chmod:
					fallthrough
				case fsnotify.Write:
					pwd.Update(event.Name)

				//deletion - remove from map
				case fsnotify.Remove:
					pwd.Remove(event.Name)
				}

				pwd.Print()
				fmt.Println(event)

				/*
					fmt.Println("event:", event)
					if event.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("modified file:", event.Name)
					}
				*/
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(pwd.path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	<-done
}
