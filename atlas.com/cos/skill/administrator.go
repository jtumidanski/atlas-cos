package skill

import "gorm.io/gorm"

type EntityFunction func() ([]string, func(e *entity))

// create a new skill for the given character and skill. Applying EntityFunction's to modify additional skill
// attributes. Returns a structure representing the skill created, or an error if one occurred.
func create(db *gorm.DB, characterId uint32, skillId uint32, modifiers ...EntityFunction) (Model, error) {
	e := &entity{
		SkillId:     skillId,
		CharacterId: characterId,
	}

	for _, modifier := range modifiers {
		_, f := modifier()
		f(e)
	}

	err := db.Create(e).Error
	if err != nil {
		return Model{}, err
	}
	return transform(*e)
}

// update a skill, applying EntityFunction's to modify the attributes of the skill. Returns an error if one occurred.
func update(db *gorm.DB, id uint32, modifiers ...EntityFunction) error {
	e := &entity{}
	var columns []string
	for _, modifier := range modifiers {
		c, u := modifier()
		columns = append(columns, c...)
		u(e)
	}
	return db.Model(&entity{ID: id}).Select(columns).Updates(e).Error
}

// setExpiration Sets the expiration of the skill.
func setExpiration(expiration int64) EntityFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Expiration"}, func(e *entity) {
			e.Expiration = expiration
		}
	}
}

// setMasterLevel Sets the master level of the skill.
func setMasterLevel(level uint32) EntityFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"MasterLevel"}, func(e *entity) {
			e.MasterLevel = level
		}
	}
}

// setLevel Sets the current level of the skill.
func setLevel(level uint32) EntityFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"SkillLevel"}, func(e *entity) {
			e.SkillLevel = level
		}
	}
}
