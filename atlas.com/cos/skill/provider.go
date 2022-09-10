package skill

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

type Transformer func(e *entity) *Model

// getById retrieves a skill by the character and skill identifiers.
func getById(characterId uint32, skillId uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{CharacterId: characterId, SkillId: skillId})
	}
}

// getForCharacter retrieves all skills for the character.
func getForCharacter(characterId uint32) database.EntityProvider[[]entity] {
	return func(db *gorm.DB) model.Provider[[]entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId})
	}
}

// transform produces a immutable skill structure
func transform(e entity) (Model, error) {
	return Model{
		id:          e.ID,
		skillId:     e.SkillId,
		level:       e.SkillLevel,
		masterLevel: e.MasterLevel,
		expiration:  e.Expiration,
	}, nil
}
