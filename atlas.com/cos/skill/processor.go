package skill

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"atlas-cos/skill/information"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ByCharacterAndIdModelProvider(db *gorm.DB) func(characterId uint32, skillId uint32) model.Provider[Model] {
	return func(characterId uint32, skillId uint32) model.Provider[Model] {
		return database.ModelProvider[Model, entity](db)(getById(characterId, skillId), transform)
	}
}

// GetSkill retrieves the identified skill for the given character.
func GetSkill(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32) (Model, error) {
	return func(characterId uint32, skillId uint32) (Model, error) {
		return ByCharacterAndIdModelProvider(db)(characterId, skillId)()
	}
}

func IfHasSkillGetEffect(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, skillId uint32) (*information.Effect, bool) {
	return func(characterId uint32, skillId uint32) (*information.Effect, bool) {
		if skill, err := GetSkill(l, db)(characterId, skillId); err == nil && skill.Level() > 0 {
			i, err := information.GetById(l, span)(skillId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve information for skill %d.", skillId)
				return nil, false
			} else {
				return &i.Effects()[skill.Level()-1], true
			}
		} else {
			return nil, false
		}
	}
}

// UpdateSkill updates the skill for the given character. Returns an error if one occurred.
func UpdateSkill(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) error {
	return func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) error {
		if skill, err := GetSkill(l, db)(characterId, skillId); err == nil {
			return update(db, skill.Id(), setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
		}
		_, err := create(db, characterId, skillId, setLevel(level), setMasterLevel(masterLevel), setExpiration(expiration))
		return err
	}
}

func ByCharacterModelProvider(db *gorm.DB) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return database.ModelSliceProvider[Model, entity](db)(getForCharacter(characterId), transform)
	}
}

// GetSkills retrieves the skills for a given character. Returns an error if one occurred.
func GetSkills(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return ByCharacterModelProvider(db)(characterId)()
	}
}

// AwardSkills awards the given character the designated skills. Returns an error if one occurred.
func AwardSkills(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, skills ...uint32) error {
	return func(characterId uint32, skills ...uint32) error {
		for _, skillId := range skills {
			err := AwardSkill(l, db, span)(characterId, skillId)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// AwardSkill awards the given character the designated skill. Returns an error if one occurred.
func AwardSkill(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, skillId uint32) error {
	return func(characterId uint32, skillId uint32) error {
		if i, err := information.GetById(l, span)(skillId); err == nil {
			maxLevel := len(i.Effects())
			_, err := create(db, characterId, skillId, setMasterLevel(uint32(maxLevel)), setExpiration(-1))
			return err
		} else {
			return errors.New("unable to locate skill information")
		}
	}
}
