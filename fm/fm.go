package fm

import (
	"sort"
	"strings"
)

//Map is unsorted, which is not as fun.
//Would need other sorting mechanisms, probably just keep this away from the actual folder object.
//not very efficient doing this every time?
func GetSortedFileList(f Folder) []string {
	list := make([]string, 0, len(f.Contents()))

	index := 0

	//https://gist.github.com/zhum/57cb45d8bbea86d87490
	for _, file := range f.Contents() {
		index = sort.Search(len(list), func(i int) bool {
			//behave more like ls
			if strings.ToLower(list[i]) == strings.ToLower(file.Name()) {
				return list[i] < file.Name()
			} else {
				return strings.ToLower(list[i]) > strings.ToLower(file.Name())
			}
		})
		list = append(list, file.Name())
		copy(list[index+1:], list[index:])
		list[index] = file.Name()
	}

	return list
}
