package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type adjustJobCommand struct {
	CharacterId uint32 `json:"characterId"`
	JobId       uint16 `json:"jobId"`
}

func AdjustJobCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &adjustJobCommand{}
	}
}

func HandleAdjustJobCommand(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*adjustJobCommand); ok {
			err := character.AdjustJob(l, db)(event.CharacterId, event.JobId)
			if err != nil {
				l.WithError(err).Errorf("Unable to adjust the job of character %d.", event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
