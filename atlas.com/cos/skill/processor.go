package skill

import (
	"atlas-cos/skill/information"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetSkill retrieves the identified skill for the given character.
func GetSkill(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32) (*Model, bool) {
	return func(characterId uint32, skillId uint32) (*Model, bool) {
		s, err := getById(db, characterId, skillId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve skill %d for character %d.", skillId, characterId)
			return nil, false
		}
		return s, true
	}
}

func IfHasSkillGetEffect(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32) (*information.Effect, bool) {
	return func(characterId uint32, skillId uint32) (*information.Effect, bool) {
		if skill, ok := GetSkill(l, db)(characterId, skillId); ok && skill.Level() > 0 {
			i, err := information.GetById(l)(skillId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve information for skill %d.", skillId)
				return nil, false
			} else {
				return &i.Effects()[skill.Level() - 1], true
			}
		} else {
			return nil, false
		}
	}
}

// UpdateSkill updates the skill for the given character. Returns an error if one occurred.
func UpdateSkill(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) error {
	return func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) error {
		if skill, ok := GetSkill(l, db)(characterId, skillId); ok {
			return update(db, skill.Id(), setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
		}
		_, err := create(db, characterId, skillId, setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
		return err
	}
}

// GetSkills retrieves the skills for a given character. Returns an error if one occurred.
func GetSkills(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		return getForCharacter(db, characterId)
	}
}

// AwardSkills awards the given character the designated skills. Returns an error if one occurred.
func AwardSkills(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skills ...uint32) error {
	return func(characterId uint32, skills ...uint32) error {
		for _, skillId := range skills {
			err := AwardSkill(l, db)(characterId, skillId)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// AwardSkill awards the given character the designated skill. Returns an error if one occurred.
func AwardSkill(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32) error {
	return func(characterId uint32, skillId uint32) error {
		if i, err := information.GetById(l)(skillId); err == nil {
			maxLevel := len(i.Effects())
			_, err := create(db, characterId, skillId, setMasterLevel(uint32(maxLevel)), setExpiration(-1))
			return err
		} else {
			return errors.New("unable to locate skill information")
		}
	}
}
