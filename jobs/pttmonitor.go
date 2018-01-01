package jobs

import (
	"net/http"
	"time"

	log "github.com/meifamily/logrus"
)

type pttMonitor struct {
	duration time.Duration
	retry    int
}

func NewPttMonitor() *pttMonitor {
	return &pttMonitor{
		duration: 1 * time.Minute,
		retry:    3,
	}
}

func (pm pttMonitor) Run() {
	log.Info("Start Ptt Monitor")

	var errorCounter = 0
	var url = "https://www.ptt.cc/bbs/index.html"
	ticker := time.NewTicker(pm.duration)
	for _ = range ticker.C {
		resp, err := http.Get(url)
		if err != nil {
			log.WithError(err).Error("HTTP Get Error")
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("Ptt is alive")
			if errorCounter >= pm.retry {
				log.Info("Ptt is back to life")
				go NewChecker().Run()
				go NewPushSumChecker().Run()
				go NewCommentChecker().Run()
			}
			errorCounter = 0
		}
		if err == nil && resp.StatusCode != http.StatusOK {
			if errorCounter < pm.retry {
				log.Info("Ptt is dying")
			}
			if errorCounter == pm.retry {
				log.Info("Ptt is Dead")
				go NewChecker().Stop()
				go NewPushSumChecker().Stop()
				go NewCommentChecker().Stop()
			}
			errorCounter++
		}
	}
}
