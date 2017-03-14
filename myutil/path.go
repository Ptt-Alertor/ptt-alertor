package myutil

import "os"

func ProjectRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func StoragePath() string {
	return ProjectRootPath() + "/storage"
}
