package fm

import (
	"time"
)

// Define an interface for file system operations that'd be nice to support.
// The way I see it, there are two types of operations in a File Manager:
// - FS-specific
//     Operations that involve actually interacting with the file system, such
//     as copying and so on.
// - "meta" operations
//	   Operations composed of these FS-specific operations, such as Undo/Redo.
//	   For example: Inverse of move is move. Inverse of copy is trash.
//     These are handled on the FM side, and so don't need to be defined per-FS.

type FileSystem interface {
	// FS navigation tools
	GetFolder(path string) (Folder, error)
	GetFile(path string) (File, error)
	// Only bother with basic features that FileInfo gives
	// (name/size/last modify date/permissions)
	// May also want to provide some sort of a "preferred order" mechanism
	// For example, it'd be useful for sorting notes if some note app implements
	// its notebook format as a filesystem.
	// Things like comments and labels could be provided on FM level.
	// Though, at the same time, they could also be provided on FS level, in
	// that case would need to define a new interface for files.
	// Need to think more about this, but for now FileInfo does the job.

	// FS operation tools

	// Copying multiple files is treated as a single operation by a FM user.
	// Doing it this way allows us to hand off the entire operation to the FS
	// instead of having the FM micromanage it. This also allows things like
	// the FS layer doing these operations in a separate process, for example.
	// This does come with the implication that conflict resolution
	// (e.g. identical file names) should take place before the operation.
	// Which is honestly how it should be to begin with.

	// With cp/mv operations, take in the source-destination map, and just go
	// ahead and complete the whole operation.
	// A source-destination relationship should exist between the files.
	Copy(files map[string]string) (<-chan IOStatus, chan<- OPStatus, error)
	Move(files map[string]string) (<-chan IOStatus, chan<- OPStatus, error)
	// With trashing operations, take in an array of paths to be trashed.
	Trash(filenames []string) (<-chan IOStatus, chan<- OPStatus, error)  // Move to FS trash folder
	Delete(filenames []string) (<-chan IOStatus, chan<- OPStatus, error) // Erase files for good

	Mkdir(name string) error

	// File open/preview tools

	// Open and preview operations should also be done on FS level to some degree
	// For remote FS this'd allow fetching remote thumbnails (as opposed to whole file)
	// Or things like opening remote files by downloading to tmp and watching changes
	// Preview mechanisms are not defined yet, neither are mechanisms for opening files.
	// xdg-open is stupid and doesn't support multiple arguments.
	// xdg-mime can tell the defaults, so can handle this on FM level on Linux.

	// Some sort of functions would also need to exist to move files across filesystems.
	Open(filenames []string) //opens file in a relevant local app (by what mechanism?)
	// Can return a string of filenames specifically on the local FS to be opened by a FM-level mechanism.

	Preview(filename string) //returns some sort of preview object for FM use
	// Preview could in theory only accept images or raw text.
	// Up to FS to provide the actual previews, in this scenario. Can also be done out-of-process.
	// Might be nice to have a mechanism for getting multiple images, lower-quality images,
	// placeholders, and so on.

	// Cross-FS operation tools

	// For moving across filesystems, it might be interesting to use readers/writers.
	// Doing this outside of FM is complicated and may not work all that well with the out-of-process IO
	// that I would like to have at least for the local FS.
	// This is not impossible, but it requires the out-of-process components be implemented on FM level.
	// That requires dealing with IPC'ing things such as remote connection credentials.
	// This is not very fun.

	// However, the idea is that these will only be used to move data between filesystems.
	// On FM level, can define a path system that deals with filesystems, then wrappers for IO functions.
	// e.g. Mkdir("/home/user/newfolder") will call the local fs Mkdir(), while
	// Mkdir("cloud://newfolder") will call the cloud fs Mkdir()
	// The FM UI layer would not care about what functions to call, and would always just use the
	// generic Copy, Move, Open, Preview, etc. functions, which will then make the necessary calls.

	//Download(filename string) //returns some Reader
	//Upload(filename string)   //returns some Writer
	// These also need to return the relevant IOStatus and OPStatus channels to pause the operation.
	// up to FM to do conflict resolution, this should always overwrite and so on, only error in case
	// of permission denial and such.
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

// Knowing currently copied file allows more granular progress reporting and displaying it in the FM.
// Knowing files created as result of an operation allows for undoing it.

// For interacting with an ongoing operation
type OPStatus int

const (
	Ongoing OPStatus = iota
	Pause
	Cancel
	Fail
	Complete
)

// conflict resolution falls to the FM.

// for FM use, folders are more "interesting" until we get to actually changing the file system
type Folder interface {
	Path() string
	Contents() map[string]File
	// Contents is a map so that it's possible to quickly add/remove content based on filename
	// This does, however, mean that you can't have multiple files with the same name in one folder.
	// A different mechanism may be good to have here, however, issue can be dealt with by always calling
	// FileInfo.Name() instead of using the map filename string.
	Watch() chan bool // ping channel if folder contents changed
	Close()           // get rid of a folder's watcher
}

// if name == "" and contents == nil, this folder object is invalid (is this the go way? idk lol)

// The idea would be to have this exist as an interface so that other fun stuff like
// virtual file systems and network fs can be added

// so if any of these get called on the UI thread I'll get sad.

// The file interface is a simpler version of the os.FileInfo interface.
// File permissions are an interesting story, and standard UNIX-style permissions are
// not flexible enough for things like cloud fs or other OS support, so I'm not sure
// what the best generic approach here would be - if there can even be one.
// Personally, I do not manage permissions via the file manager all that often, so
// I could just leave it out completely.
type File interface {
	Name() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}

// For file permissions and per-FS settings, I could provide some simple building blocks
// that could allow file systems to provide simple UIs for managing those - I will need
// to have this exist in some form either way for handling things like authentication,
// for example. Cover a few basics - enough to flip some switches for the FS config and
// basic permission handling. Could also straight up allow web views for authentication
// and dealing with the complicated world of cloud sharing.

// Remember, this is all for a GUI-based file manager platform.
