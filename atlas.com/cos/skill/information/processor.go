package information

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"github.com/sirupsen/logrus"
)

func GetSkillInformation(l logrus.FieldLogger) func(skillId uint32) (*Model, bool) {
	return func(skillId uint32) (*Model, bool) {
		s, err := requests.Skill().GetById(skillId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve skill %d information.", skillId)
			return nil, false
		}
		return makeSkill(s.Data), true
	}
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
