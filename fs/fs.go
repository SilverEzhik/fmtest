package fs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Mutex it
type folder struct {
	path     string
	contents map[string]os.FileInfo
}

//get the folder struct
func GetFolder(path string) *folder {
	f := new(folder)
	*f = folder{path: path}
	f.contents = make(map[string]os.FileInfo)
	f.refresh()

	return f
}

func (f *folder) Path() string {
	return f.path
}
func (f *folder) Contents() map[string]os.FileInfo {
	return f.contents
}

//notifies the given channel about changes to the folder
func (f *folder) Watch(hasUpdated chan bool) {
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
					f.updateItem(event.Name)

				//deletion - remove from map
				case fsnotify.Remove:
					f.removeItem(event.Name)
				}

				hasUpdated <- true

			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(f.path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	<-done
}

// Get initial folder contents from file system
// Might be nice to have a mechanism that would just go ahead and refresh if
// Many update/remove events are queued.
func (f *folder) refresh() {
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
