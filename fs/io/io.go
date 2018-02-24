package fs

import "os/exec"

func CopyFile(source, destination string) error {
	cmd := exec.Command("cp", source, destination)
	_, err := cmd.Output()
	if err != nil {
		return err
	} else {
		return nil
	}
}
