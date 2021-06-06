package jobs

import (
	"sync"

	log "github.com/meifamily/logrus"

	"time"

	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/board"
)

type Fetcher struct {
}

func NewFetcher() *Fetcher {
	f := new(Fetcher)
	return f
}

func (f Fetcher) Run() {
	boards := models.Board().All()

	var wg sync.WaitGroup
	for _, bd := range boards {
		wg.Add(1)
		go func(bd board.Board) {
			defer wg.Done()
			bd.Articles = bd.FetchArticles()
			bd.Save()
			log.WithField("board", bd.Name).Info("Fetched")
		}(*bd)
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
	log.Info("All fetcher done")
}
