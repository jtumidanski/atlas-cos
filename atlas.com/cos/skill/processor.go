package skill

import (
	"atlas-cos/skill/information"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

// GetSkill retrieves the identified skill for the given character.
func (p processor) GetSkill(characterId uint32, skillId uint32) (*Model, bool) {
	s, err := getById(p.db, characterId, skillId)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve skill %d for character %d.", skillId, characterId)
		return nil, false
	}
	return s, true
}

// UpdateSkill updates the skill for the given character. Returns an error if one occurred.
func (p processor) UpdateSkill(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) error {
	if skill, ok := p.GetSkill(characterId, skillId); ok {
		return update(p.db, skill.Id(), setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
	}
	_, err := create(p.db, characterId, skillId, setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
	return err
}

// GetSkills retrieves the skills for a given character. Returns an error if one occurred.
func (p processor) GetSkills(characterId uint32) ([]*Model, error) {
	return getForCharacter(p.db, characterId)
}

// AwardSkills awards the given character the designated skills. Returns an error if one occurred.
func (p processor) AwardSkills(characterId uint32, skills ...uint32) error {
	for _, skillId := range skills {
		err := p.AwardSkill(characterId, skillId)
		if err != nil {
			return err
		}
	}
	return nil
}

// AwardSkill awards the given character the designated skill. Returns an error if one occurred.
func (p processor) AwardSkill(characterId uint32, skillId uint32) error {
	if i, ok := information.Processor(p.l, p.db).GetSkillInformation(skillId); ok {
		maxLevel := len(i.Effects())
		_, err := create(p.db, characterId, skillId, setMasterLevel(uint32(maxLevel)), setExpiration(-1))
		return err
	} else {
		return errors.New("unable to locate skill information")
	}
}
