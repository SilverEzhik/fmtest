package fs

import (
	"os"
	"syscall"
)

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
