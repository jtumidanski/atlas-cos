package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
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
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*resetAPCommand); ok {
			err := character.ResetAP(l, db)(event.CharacterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to reset AP of character %d.", event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
