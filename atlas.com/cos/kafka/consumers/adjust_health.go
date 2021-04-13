package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type adjustHealthCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint16 `json:"amount"`
}

func AdjustHealthCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &adjustHealthCommand{}
	}
}

func HandleAdjustHealthCommand(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*adjustHealthCommand); ok {
			character.Processor(l, db).AdjustHealth(event.CharacterId, event.Amount)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [AdjustHealthCommand]")
		}
	}
}