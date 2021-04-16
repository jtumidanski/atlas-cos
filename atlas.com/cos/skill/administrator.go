package skill

import "gorm.io/gorm"

type EntityUpdateFunction func(e entity)

func Create(db *gorm.DB, characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration uint64) (*Model, error) {
	e := &entity{
		SkillId:     skillId,
		CharacterId: characterId,
		SkillLevel:  level,
		MasterLevel: masterLevel,
		Expiration:  expiration,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeSkill(e), nil
}

func Update(db *gorm.DB, id uint32, modifiers ...EntityUpdateFunction) error {
	c := entity{ID: id}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	for _, modifier := range modifiers {
		modifier(c)
	}

	err = db.Save(&c).Error
	return err
}

func SetExpiration(expiration uint64) EntityUpdateFunction {
	return func(e entity) {
		e.Expiration = expiration
	}
}

func SetMasterLevel(level uint32) EntityUpdateFunction {
	return func(e entity) {
		e.MasterLevel = level
	}
}

func SetLevel(level uint32) EntityUpdateFunction {
	return func(e entity) {
		e.SkillLevel = level
	}
}
