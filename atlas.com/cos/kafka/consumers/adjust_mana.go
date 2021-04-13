package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type adjustManaCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint16 `json:"amount"`
}

func AdjustManaCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &adjustManaCommand{}
	}
}

func HandleAdjustManaCommand(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*adjustManaCommand); ok {
			character.Processor(l, db).AdjustMana(event.CharacterId, event.Amount)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [AdjustManaCommand]")
		}
	}
}