package main

import (
	"os"
	//"time"
)

// for FM use, folders are more "interesting" until we get to actually changing the file system
type Folder interface {
	Path() string
	Contents() map[string]os.FileInfo
	Watch() chan bool // ping channel if folder contents changed
	Close()           // get rid of a folder's watcher
}

// if name == "" and contents == nil, this folder object is invalid (is this the go way? idk lol)

// The idea would be to have this exist as an interface so that other fun stuff like
// virtual file systems and network can be added

// so if any of these get called on the UI thread I'll get sad.
