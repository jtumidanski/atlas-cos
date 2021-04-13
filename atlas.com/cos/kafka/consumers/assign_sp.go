package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type assignSPCommand struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func AssignSPCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &assignSPCommand{}
	}
}

func HandleAssignSPCommand(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*assignSPCommand); ok {
			character.Processor(l, db).AssignSP(event.CharacterId, event.SkillId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [AssignSPCommand]")
		}
	}
}
