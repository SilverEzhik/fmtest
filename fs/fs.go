package fs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Mutex this
type folder struct {
	path     string
	contents map[string]os.FileInfo
	watcher  chan bool
	done     chan struct{}
}

//get the folder struct
func GetFolder(path string) *folder {
	f := &folder{path: path, contents: make(map[string]os.FileInfo)}
	f.Refresh()

	return f
}

// Get folder path
func (f *folder) Path() string {
	return f.path
}

// Get folder contents
func (f *folder) Contents() map[string]os.FileInfo {
	return f.contents
}

// Get channel for notifications on changes to the folder
func (f *folder) Watch() chan bool {
	f.Refresh()
	f.done = make(chan struct{})
	f.watcher = make(chan bool)
	go f.fsWatcher()
	return f.watcher
}

func (f *folder) Close() {
	close(f.done)
}

func (f *folder) fsWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	defer watcher.Close()
	defer close(f.watcher)

	err = watcher.Add(f.path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

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
				f.updateItem(event.Name)

			//deletion - remove from map
			case fsnotify.Remove:
				f.removeItem(event.Name)
			}

			fmt.Println(f.path, "(w) -", event.Name, event.Op)
			f.watcher <- true
		case err := <-watcher.Errors:
			fmt.Println("error:", err)
		case <-f.done:
			return
		}
	}
}

// Get initial folder contents from file system
// Might be nice to have a mechanism that would just go ahead and refresh if
// Many update/remove events are queued.
func (f *folder) Refresh() {
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
func (f *folder) updateItem(path string) {
	file, err := os.Stat(path)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	f.contents[filepath.Base(path)] = file
	//fmt.Printf("up: %s\n", f.contents[filepath.Base(path)].Name())
}

func (f *folder) removeItem(path string) {
	delete(f.contents, filepath.Base(path))
}

// moving files is tricky. fsnotify does not link "rename" and "create" events, which means that it is possible to misinterpret these events in a situation where a lot of events happen in a folder.
// what do?
