package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
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
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*gainLevelEvent); ok {
			character.Processor(l, db).GainLevel(event.CharacterId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [GainLevelEvent]")
		}
	}
}
