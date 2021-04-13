package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type adjustMesoCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint32 `json:"amount"`
	Show        bool   `json:"show"`
}

func AdjustMesoCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &adjustMesoCommand{}
	}
}

func HandleAdjustMesoCommand(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*adjustMesoCommand); ok {
			character.Processor(l, db).AdjustMeso(event.CharacterId, event.Amount, event.Show)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [AdjustMesoCommand]")
		}
	}
}
