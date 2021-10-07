package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type adjustHealthCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func AdjustHealthCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &adjustHealthCommand{}
	}
}

func HandleAdjustHealthCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*adjustHealthCommand); ok {
			l.Debugf("Begin event handling.")
			if event.Amount == 0 {
				l.Infoln("Received erroneous command to adjust by 0. This should be cleaned up.")
			} else {
				character.AdjustHealth(l, db, span)(event.CharacterId, event.Amount)
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
