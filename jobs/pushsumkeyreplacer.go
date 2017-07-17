package jobs

import (
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models/pushsum"
)

type PushSumKeyReplacer struct{}

func NewPushSumKeyReplacer() *PushSumKeyReplacer {
	return &PushSumKeyReplacer{}
}

func (r PushSumKeyReplacer) Run() {
	err := pushsum.ReplaceBaseKeys()
	if err != nil {
		log.WithError(err).Error("Replace Pushsum Key Failed")
	}
	log.Info("Replace Pushsum Key Done")
}
