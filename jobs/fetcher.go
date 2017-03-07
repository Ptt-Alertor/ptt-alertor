package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"sync"

	board "github.com/liam-lai/ptt-alertor/models/ptt/board/file"
	"github.com/liam-lai/ptt-alertor/myutil"
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
	for _, bd := range boards {
		wg.Add(1)
		go func(bd board.Board) {
			defer wg.Done()
			fmt.Println(bd.Name)
			articles := bd.OnlineArticles()
			articlesJSON, err := json.Marshal(articles)
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(f.workingDir+bd.Name+".json", articlesJSON, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}(bd)
	}
	wg.Wait()
	log.Println("fetcher done")
}
