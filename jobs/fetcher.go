package jobs

import (
	"fmt"
	"log"
	"sync"

	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
)

type Fetcher struct {
}

func NewFetcher() *Fetcher {
	f := new(Fetcher)
	return f
}

func (f Fetcher) Run() {
	boards := new(board.Board).All()

	var wg sync.WaitGroup
	for _, bd := range boards {
		wg.Add(1)
		go func(bd board.Board) {
			defer wg.Done()
			fmt.Println(bd.Name)
			bd.Articles = bd.OnlineArticles()
			bd.Save()
		}(*bd)
	}
	wg.Wait()
	log.Println("fetcher done")
}
