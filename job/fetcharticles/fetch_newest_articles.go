package main

import (
	"io/ioutil"
	"log"

	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

func main() {
	var boards []string

	articlesDir := myutil.StoragePath() + "/articles/"
	files, _ := ioutil.ReadDir(articlesDir)
	for _, file := range files {
		boardName, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		boards = append(boards, boardName)
	}

	for _, board := range boards {
		articlesJSON := pttboard.Index(board)
		err := ioutil.WriteFile(myutil.ProjectRootPath()+"/storage/articles/"+board+".json", articlesJSON, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
