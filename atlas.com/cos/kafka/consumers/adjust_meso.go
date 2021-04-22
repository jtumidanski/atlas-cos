package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type adjustMesoCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int32  `json:"amount"`
	Show        bool   `json:"show"`
}

func AdjustMesoCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &adjustMesoCommand{}
	}
}

func HandleAdjustMesoCommand(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*adjustMesoCommand); ok {
			l.Debugf("Begin event handling.")
			character.Processor(l, db).AdjustMeso(event.CharacterId, event.Amount, event.Show)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
