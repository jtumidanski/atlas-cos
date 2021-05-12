package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type assignSPCommand struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func AssignSPCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &assignSPCommand{}
	}
}

func HandleAssignSPCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*assignSPCommand); ok {
			l.Debugf("Begin event handling.")
			character.AssignSP(l, db)(event.CharacterId, event.SkillId)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
