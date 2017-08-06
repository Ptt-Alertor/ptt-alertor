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
		}
		if err == nil && resp.StatusCode != http.StatusOK {
			if errorCounter < 3 {
				log.Info("Ptt is dying")
			}
			errorCounter++
		}
		if err == nil && resp.StatusCode == http.StatusOK && errorCounter != 0 {
			log.Info("Ptt is back to life")
			errorCounter = 0
			go NewChecker().Run()
			go NewPushSumChecker().Run()
			go NewPushListChecker().Run()
		}
		if errorCounter == 3 {
			log.Info("Ptt is Dead")
			go NewChecker().Stop()
			go NewPushSumChecker().Stop()
			go NewPushListChecker().Stop()
		}
		time.Sleep(pm.duration)
	}
}
