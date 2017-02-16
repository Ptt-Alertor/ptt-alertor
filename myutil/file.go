package myutil

import "strings"

func FileNameAndExtension(basefilename string) (string, string) {
	filenames := strings.Split(basefilename, ".")
	extension := filenames[len(filenames)-1]
	filename := strings.Join(filenames[:len(filenames)-1], ".")
	return filename, extension
}
