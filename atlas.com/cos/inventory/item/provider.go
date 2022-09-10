package item

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

func getById(id uint32) database.EntityProvider[entityInventoryItem] {
	return func(db *gorm.DB) model.Provider[entityInventoryItem] {
		return database.Query[entityInventoryItem](db, &entityInventoryItem{ID: id})
	}
}

func getForCharacter(inventoryId uint32, itemId uint32) database.EntityProvider[[]entityInventoryItem] {
	return func(db *gorm.DB) model.Provider[[]entityInventoryItem] {
		return database.SliceQuery[entityInventoryItem](db, &entityInventoryItem{InventoryId: inventoryId, ItemId: itemId})
	}
}

func getBySlot(inventoryId uint32, slot int16) database.EntityProvider[entityInventoryItem] {
	return func(db *gorm.DB) model.Provider[entityInventoryItem] {
		return database.Query[entityInventoryItem](db, &entityInventoryItem{InventoryId: inventoryId, Slot: slot})
	}
}

func getByInventory(inventoryId uint32) database.EntityProvider[[]entityInventoryItem] {
	return func(db *gorm.DB) model.Provider[[]entityInventoryItem] {
		return database.SliceQuery[entityInventoryItem](db, &entityInventoryItem{InventoryId: inventoryId})
	}
}

func getByEquipmentId(equipmentId uint32) database.EntityProvider[entityInventoryItem] {
	return func(db *gorm.DB) model.Provider[entityInventoryItem] {
		return database.Query[entityInventoryItem](db, &entityInventoryItem{ReferenceId: equipmentId})
	}
}

func getItemAttributes(id uint32) database.EntityProvider[entityItem] {
	return func(db *gorm.DB) model.Provider[entityItem] {
		return database.Query[entityItem](db, &entityItem{ID: id})
	}
}

func minFreeSlot(items []Model) int16 {
	slot := int16(1)
	i := 0

	for {
		if i >= len(items) {
			return slot
		} else if slot < items[i].Slot() {
			return slot
		} else if slot == items[i].Slot() {
			slot += 1
			i += 1
		} else if items[i].Slot() <= 0 {
			i += 1
		}
	}
}
