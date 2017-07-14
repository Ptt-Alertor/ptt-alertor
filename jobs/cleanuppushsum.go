package jobs

import (
	"strconv"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/pushsum"
)

type cleanUpPushSum struct {
}

func NewCleanUpPushSum() *cleanUpPushSum {
	return &cleanUpPushSum{}
}

func (c cleanUpPushSum) Run() {
	today := time.Now().Day()
	err := pushsum.DelDayKeys(strconv.Itoa(today))
	if err != nil {
		log.WithError(err).Error("Clean Up Today Keys Failed")
	}
	yesterday := time.Now().AddDate(0, 0, -1).Day()
	err = pushsum.DelDayKeys(strconv.Itoa(yesterday))
	if err != nil {
		log.WithError(err).Error("Clean Up Overdue Keys Failed")
	}
	log.Info("Clean Up Overdue Keys Done")
	err = pushsum.ReplacePreKeys()
	if err != nil {
		log.WithError(err).Error("Clean Up Overdue Keys Failed")
	}
	log.Info("Change Pre Key to Base Key Done")
}
