package inventory

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

func get(characterId uint32, inventoryType int8) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{CharacterId: characterId, InventoryType: inventoryType})
	}
}
