package skill

import "gorm.io/gorm"

func GetById(db *gorm.DB, characterId uint32, skillId uint32) (*Model, error) {
	var result entity
	err := db.Where(&entity{CharacterId: characterId, SkillId: skillId}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return makeSkill(&result), nil
}

func GetForCharacter(db *gorm.DB, characterId uint32) ([]*Model, error) {
	var result []entity
	err := db.Where(&entity{CharacterId: characterId}).Find(&result).Error
	if err != nil {
		return nil, err
	}

	var skills = make([]*Model, 0)
	for _, r := range result {
		skills = append(skills, makeSkill(&r))
	}
	return skills, nil
}

func makeSkill(e *entity) *Model {
	return &Model{
		id:          e.ID,
		level:       e.SkillLevel,
		masterLevel: e.MasterLevel,
		expiration:  e.Expiration,
	}
}
