package jobs

import (
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models/pushsum"
)

type ReplacePushSumKey struct{}

func NewReplacePushSumKey() *ReplacePushSumKey {
	return &ReplacePushSumKey{}
}

func (r ReplacePushSumKey) Run() {
	err := pushsum.ReplaceBaseKeys()
	if err != nil {
		log.WithError(err).Error("Replace Pushsum Key Failed")
	}
	log.Info("Replace Pushsum Key Done")
}
