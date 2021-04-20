package consumers

import (
	"atlas-cos/character"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type assignAPCommand struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func AssignAPCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &assignAPCommand{}
	}
}

func HandleAssignAPCommand(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*assignAPCommand); ok {
			l.Debugf("Begin event handling.")
			switch event.Type {
			case "STRENGTH":
				character.Processor(l, db).AssignStrength(event.CharacterId)
				break
			case "DEXTERITY":
				character.Processor(l, db).AssignDexterity(event.CharacterId)
				break
			case "INTELLIGENCE":
				character.Processor(l, db).AssignIntelligence(event.CharacterId)
				break
			case "LUCK":
				character.Processor(l, db).AssignLuck(event.CharacterId)
				break
			case "HP":
				character.Processor(l, db).AssignHp(event.CharacterId)
				break
			case "MP":
				character.Processor(l, db).AssignMp(event.CharacterId)
				break
			}
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
