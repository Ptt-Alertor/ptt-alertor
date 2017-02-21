package main

import (
	"io/ioutil"
	"log"

	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/ptt/board"
)

func main() {
	var boards []board.Board

	articlesDir := myutil.StoragePath() + "/articles/"
	files, _ := ioutil.ReadDir(articlesDir)
	for _, file := range files {
		boardName, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		b := new(board.Board)
		b.Name = boardName
		boards = append(boards, *b)
	}

	for _, board := range boards {
		articlesJSON := board.IndexJSON()
		err := ioutil.WriteFile(articlesDir+board.Name+".json", articlesJSON, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
