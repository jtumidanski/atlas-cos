package skill

import "gorm.io/gorm"

type Transformer func(e *entity) *Model

// getOne queries the database for the first entity which meets the provided query criteria, then applies the
// Transformer function to produce a result.
func getOne(db *gorm.DB, query interface{}, transformer Transformer) (*Model, error) {
	var result = &entity{}
	err := db.Where(query).First(result).Error
	if err != nil {
		return nil, err
	}
	return transformer(result), nil
}

// getArray queries the database for the entities which meets the provided query criteria, then applies the
// Transformer function to produce a array of results.
func getArray(db *gorm.DB, query interface{}, transformer Transformer) ([]*Model, error) {
	var result []entity
	err := db.Where(query).Find(&result).Error
	if err != nil {
		return nil, err
	}

	var skills = make([]*Model, 0)
	for _, r := range result {
		skills = append(skills, transformer(&r))
	}
	return skills, nil
}

// getById retrieves a skill by the character and skill identifiers.
func getById(db *gorm.DB, characterId uint32, skillId uint32) (*Model, error) {
	return getOne(db, &entity{CharacterId: characterId, SkillId: skillId}, transform)
}

// getForCharacter retrieves all skills for the character.
func getForCharacter(db *gorm.DB, characterId uint32) ([]*Model, error) {
	return getArray(db, &entity{CharacterId: characterId}, transform)
}

// Produces a immutable skill structure
func transform(e *entity) *Model {
	return &Model{
		id:          e.ID,
		level:       e.SkillLevel,
		masterLevel: e.MasterLevel,
		expiration:  e.Expiration,
	}
}
