package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
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
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*assignAPCommand); ok {
			l.Debugf("Begin event handling.")
			switch event.Type {
			case "STRENGTH":
				character.AssignStrength(l, db)(event.CharacterId)
				break
			case "DEXTERITY":
				character.AssignDexterity(l, db)(event.CharacterId)
				break
			case "INTELLIGENCE":
				character.AssignIntelligence(l, db)(event.CharacterId)
				break
			case "LUCK":
				character.AssignLuck(l, db)(event.CharacterId)
				break
			case "HP":
				character.AssignHp(l, db)(event.CharacterId)
				break
			case "MP":
				character.AssignMp(l, db)(event.CharacterId)
				break
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
