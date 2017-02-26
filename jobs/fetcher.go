package jobs

import (
	"fmt"
	"io/ioutil"
	"log"

	"sync"

	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/ptt/board"
)

type Fetcher struct {
	workingDir string
}

func NewFetcher() *Fetcher {
	f := new(Fetcher)
	f.workingDir = myutil.StoragePath() + "/articles/"
	return f
}

func (f Fetcher) Run() {
	var boards []board.Board

	files, _ := ioutil.ReadDir(f.workingDir)
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
			err := ioutil.WriteFile(f.workingDir+b.Name+".json", articlesJSON, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}(b)
	}
	wg.Wait()
	fmt.Println("fetcher done")
}
