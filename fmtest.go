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

type Folder struct {
	path     string
	contents map[string]os.FileInfo
}

// Get initial folder contents from file system
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

//takes an absolute path to file and stats it
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
	f.GetSortedFileList()
}

func (f *Folder) GetSortedFileList() []string {
	list := make([]string, 0, len(f.contents))

	index := 0

	for name := range f.contents {
		index = sort.Search(len(list), func(i int) bool {
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

	fmt.Print("files: ")
	for _, file := range list {
		fmt.Print(file, " ")
	}
	fmt.Print("\n")

	return list
}

//func watchFolder(f *Folder) {

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
