package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"sync"

	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/ptt/board"
)

var articlesDir string = myutil.StoragePath() + "/articles/"

func main() {
	var boards []board.Board

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

	var wg sync.WaitGroup
	for _, b := range boards {
		wg.Add(1)
		go func(b board.Board) {
			defer wg.Done()
			fmt.Println(b.Name)
			articlesJSON := b.IndexJSON()
			err := ioutil.WriteFile(articlesDir+b.Name+".json", articlesJSON, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}(b)
	}
	wg.Wait()
}
