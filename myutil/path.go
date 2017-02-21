package myutil

import "os"

func ProjectRootPath() string {
	dir, _ := os.Getwd()
	return dir
}

func StoragePath() string {
	return ProjectRootPath() + "/storage"
}
