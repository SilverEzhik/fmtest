package fs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// something to think about: see if it's possible to have some structure here that could help avoid
// creating multiple watchers for a single folder

type folder struct {
	path     string
	contents map[string]os.FileInfo
	mutex    sync.Mutex
	watcher  chan bool
	done     chan struct{}
	uid      uint64
}

//get the folder struct
func GetFolder(path string) (*folder, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f := &folder{path: absPath}
	f.Refresh()
	f.done = make(chan struct{})
	f.watcher = make(chan bool)
	go f.fsWatcher()

	return f, nil
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

	if f.path != filepath.Dir(f.path) {
		err = watcher.Add(filepath.Dir(f.path))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
	}

	folderRenamed := false

	for {
		select {
		case event := <-watcher.Events:
			//since we are watching the parent as well (most likely), check path.
			if filepath.Dir(event.Name) == f.path {
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
				f.watcher <- true
			} else if event.Name == f.path && event.Op == fsnotify.Rename {
				folderRenamed = true

				//this is... quite a thing to do.
				go func() {
					time.Sleep(10 * time.Millisecond)
					if folderRenamed == true {
						fmt.Println("folder gone")
						f.cleanup()
					}
				}()
			} else if folderRenamed == true && event.Op == fsnotify.Create {
				//compare folder uids
				if f.uid != getPathUID(event.Name) {
					f.cleanup()
					continue
				}

				//if this is indeed the folder we were in, change this folder object
				fmt.Println("caught new folder")
				absPath, err := filepath.Abs(event.Name)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				//stop watching old path
				err = watcher.Remove(f.path)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				f.path = absPath
				f.Refresh()
				//start watching new path
				err = watcher.Add(f.path)
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				folderRenamed = false
			}

			fmt.Println(f.path, "(w) -", event.Name, event.Op)
		case err := <-watcher.Errors:
			fmt.Println("error:", err)
		case <-f.done:
			fmt.Println("watcher over")
			return
		}
	}
}

func (f *folder) cleanup() {
	if f.path == "" && f.contents == nil {
		return
	}
	f.path = ""
	f.contents = nil
	f.Close()
}

// Get initial folder contents from file system
// Might be nice to have a mechanism that would just go ahead and refresh if
// many update/remove events are queued.
func (f *folder) Refresh() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	// replace the map
	f.contents = make(map[string]os.FileInfo)

	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f.contents[file.Name()] = file
		//fmt.Println(file.Name())
	}

	f.uid = getPathUID(f.path)
}

//unix only
//figure out a drop-in uid function for other os?
func getPathUID(path string) uint64 {
	fileinfo, _ := os.Stat(path)
	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		// 0 in inodes indicates an error, so...
		return 0
	}

	return stat.Ino
}

// Takes an absolute path to file and stats it
// Also functions as an add function
func (f *folder) updateItem(path string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	file, err := os.Stat(path)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	f.contents[filepath.Base(path)] = file
	//fmt.Printf("up: %s\n", f.contents[filepath.Base(path)].Name())
}

func (f *folder) removeItem(path string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	delete(f.contents, filepath.Base(path))
}

// moving files is tricky. fsnotify does not link "rename" and "create" events, which means that it is possible to misinterpret these events in a situation where a lot of events happen in a folder.
// what do?
