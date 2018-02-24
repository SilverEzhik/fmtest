package fm

import (
	"time"
)

type FileSystem interface {
	GetFolder(path string) (Folder, error)
	Stat(path string) (File, error)

	Copy(files map[string]string) (<-chan IOStatus, chan<- OPStatus, error)
	Move(files map[string]string) (<-chan IOStatus, chan<- OPStatus, error)
	//Copy(files []string, destination string) (<-chan IOStatus, chan<- OPStatus, error)
	//Move(files []string, destination string) (<-chan IOStatus, chan<- OPStatus, error)

	Trash(filenames []string) (<-chan IOStatus, chan<- OPStatus, error)  // Move to FS trash folder
	Delete(filenames []string) (<-chan IOStatus, chan<- OPStatus, error) // Erase files for good

	Mkdir(name string) error

	Open(filenames []string) //opens file in a relevant local app (by what mechanism?)

	Preview(filename string) //returns some sort of preview object for FM use

	//Download(filename string) //returns some Reader
	//Upload(filename string)   //returns some Writer
}

// progress tracking
type IOStatus struct {
	CurrentFile  string            // get path of file currently worked on
	FileProgress int               // progress of this individual file %
	Progress     int               // total progress %
	Results      map[string]string // get paths of whatever files emerged as a result of completing the operation
	Status       OPStatus          // current status of the operation
	Error        error
}

type OPStatus int

const (
	Ongoing OPStatus = iota
	Pause
	Cancel
	Fail
	Complete
)

type Folder interface {
	Path() string
	Contents() map[string]File
	Watch() chan bool // ping channel if folder contents changed
	Close()           // get rid of a folder's watcher
}

type File interface {
	Name() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}
