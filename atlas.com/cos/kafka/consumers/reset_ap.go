package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type resetAPCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func ResetAPCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &resetAPCommand{}
	}
}

func HandleResetAPCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*resetAPCommand); ok {
			err := character.ResetAP(l, db, span)(event.CharacterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to reset AP of character %d.", event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
