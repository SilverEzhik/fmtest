package fs

import (
	"os"
	"reflect"
)

// as seen in https://github.com/hartfordfive/protologbeat/blob/master/vendor/github.com/elastic/beats/filebeat/input/file/file_windows.go
func getPathUID(path string) uint64 {
	fileinfo, _ := os.Stat(path)

	os.SameFile(fileinfo, fileinfo)

	fileStat := reflect.ValueOf(fileinfo).Elem()
	uid := uint64(fileStat.FieldByName("idxhi").Uint()) << 32
	uid += uint64(fileStat.FieldByName("idxlo").Uint())

	return uid
}
