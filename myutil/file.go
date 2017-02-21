package myutil

import (
	"os"
	"strings"
)

func FileNameAndExtension(basefilename string) (string, string) {
	filenames := strings.Split(basefilename, ".")
	extension := filenames[len(filenames)-1]
	filename := strings.Join(filenames[:len(filenames)-1], ".")
	return filename, extension
}

func JsonFile(file os.FileInfo) (string, bool) {
	if file.IsDir() {
		return "", false
	}

	fileName, extension := FileNameAndExtension(file.Name())
	if extension != "json" {
		return "", false
	}
	return fileName, true
}
