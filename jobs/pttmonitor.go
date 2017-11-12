package jobs

import (
	"net/http"
	"time"

	log "github.com/meifamily/logrus"
)

type pttMonitor struct {
	duration time.Duration
}

func NewPttMonitor() *pttMonitor {
	return &pttMonitor{
		duration: 1 * time.Minute,
	}
}

func (pm pttMonitor) Run() {
	log.Info("Start Ptt Monitor")

	var errorCounter = 0
	var url = "https://www.ptt.cc/index.html"
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.WithError(err).Error("HTTP Get Error")
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("Ptt is alive")
			if errorCounter >= 3 {
				log.Info("Ptt is back to life")
				go NewChecker().Run()
				go NewPushSumChecker().Run()
				go NewCommentChecker().Run()
			}
			errorCounter = 0
		}
		if err == nil && resp.StatusCode != http.StatusOK {
			if errorCounter < 3 {
				log.Info("Ptt is dying")
			}
			if errorCounter == 3 {
				log.Info("Ptt is Dead")
				go NewChecker().Stop()
				go NewPushSumChecker().Stop()
				go NewCommentChecker().Stop()
			}
			errorCounter++
		}
		time.Sleep(pm.duration)
	}
}
