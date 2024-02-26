package vm

import (
	"bufio"
	"os"
	"strings"
)

func GetDisplayNameFromPath(path string) (string, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	s := bufio.NewScanner(fd)
	for s.Scan() {
		if strings.Contains(s.Text(), "displayName") {
			return strings.TrimSpace(strings.Split(s.Text(), "=")[1]), nil
		}
	}
	return "", nil
}
