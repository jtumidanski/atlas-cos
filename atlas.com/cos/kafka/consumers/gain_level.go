package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type gainLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func GainLevelEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &gainLevelEvent{}
	}
}

func HandleGainLevelEvent(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*gainLevelEvent); ok {
			l.Debugf("Begin event handling.")
			character.Processor(l, db).GainLevel(event.CharacterId)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorln("Unable to cast event provided to handler")
		}
	}
}
