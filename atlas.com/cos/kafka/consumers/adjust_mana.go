package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type adjustManaCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func AdjustManaCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &adjustManaCommand{}
	}
}

func HandleAdjustManaCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*adjustManaCommand); ok {
			l.Debugf("Begin event handling.")
			if event.Amount == 0 {
				l.Infoln("Received erroneous command to adjust by 0. This should be cleaned up.")
			} else {
				character.AdjustMana(l, db, span)(event.CharacterId, event.Amount)
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
