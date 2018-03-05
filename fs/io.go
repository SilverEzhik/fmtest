package fs

import (
	"io"
	"log"
	"os"
)

func CopyFile(source, destination string) error {
	from, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	//file permissions
	var perm os.FileMode
	stat, err := os.Stat(source)
	if err != nil {
		perm = 0666
	} else {
		perm = stat.Mode()
	}

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
