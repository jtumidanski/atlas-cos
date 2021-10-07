package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type assignAPCommand struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func AssignAPCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &assignAPCommand{}
	}
}

func HandleAssignAPCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*assignAPCommand); ok {
			l.Debugf("Begin event handling.")
			switch event.Type {
			case "STRENGTH":
				character.AssignStrength(l, db, span)(event.CharacterId)
				break
			case "DEXTERITY":
				character.AssignDexterity(l, db, span)(event.CharacterId)
				break
			case "INTELLIGENCE":
				character.AssignIntelligence(l, db, span)(event.CharacterId)
				break
			case "LUCK":
				character.AssignLuck(l, db, span)(event.CharacterId)
				break
			case "HP":
				character.AssignHp(l, db, span)(event.CharacterId)
				break
			case "MP":
				character.AssignMp(l, db, span)(event.CharacterId)
				break
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
