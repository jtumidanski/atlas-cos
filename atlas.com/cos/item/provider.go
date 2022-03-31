package item

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

func getById(id uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{Id: id})
	}
}

func getItemsForCharacter(characterId uint32, inventoryType int8, itemId uint32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId, InventoryType: inventoryType, ItemId: itemId})
	}
}

func getItemForCharacter(characterId uint32, inventoryType int8, slot int16) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{CharacterId: characterId, InventoryType: inventoryType, Slot: slot})
	}
}

func getItemsForCharacterByInventory(characterId uint32, inventoryType int8) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId, InventoryType: inventoryType})
	}
}

func makeItem(e entity) (Model, error) {
	return Model{
		id:            e.Id,
		characterId:   e.CharacterId,
		inventoryType: e.InventoryType,
		itemId:        e.ItemId,
		quantity:      e.Quantity,
		slot:          e.Slot,
	}, nil
}
