package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type gainLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func GainLevelEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &gainLevelEvent{}
	}
}

func HandleGainLevelEvent(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*gainLevelEvent); ok {
			l.Debugf("Begin event handling.")
			character.GainLevel(l, db)(event.CharacterId)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorln("Unable to cast event provided to handler")
		}
	}
}
