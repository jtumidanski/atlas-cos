package information

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetSkillInformation(skillId uint32) (*Model, bool) {
	s, err := requests.Skill().GetById(skillId)
	if err != nil {
		p.l.Printf("[ERROR] unable to retrieve skill %d information.", skillId)
		return nil, false
	}
	return makeSkill(s.Data), true
}

func makeSkill(data attributes.SkillData) *Model {
	var effects = make([]Effect, 0)
	for range data.Attributes.Effects {
		effects = append(effects, Effect{})
	}

	return &Model{
		effects: effects,
	}
}
