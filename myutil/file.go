package myutil

import (
	"os"
	"path/filepath"
	"strings"
)

func FileNameAndExtension(basefilename string) (filename string, ext string) {
	extWithDot := filepath.Ext(basefilename)
	filename = strings.TrimRight(basefilename, extWithDot)
	ext = strings.TrimPrefix(extWithDot, ".")
	return filename, ext
}

func JSONFile(file os.FileInfo) (fileName string, ok bool) {
	if file.IsDir() {
		return "", false
	}

	fileName, ext := FileNameAndExtension(file.Name())
	if ext != "json" {
		return "", false
	}
	return fileName, true
}
