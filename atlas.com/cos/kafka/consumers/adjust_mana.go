package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*adjustManaCommand); ok {
			l.Debugf("Begin event handling.")
			character.Processor(l, db).AdjustMana(event.CharacterId, event.Amount)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
