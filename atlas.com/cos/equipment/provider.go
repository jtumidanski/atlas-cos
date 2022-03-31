package equipment

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

func getByEquipmentId(equipmentId uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{EquipmentId: equipmentId})
	}
}

func getById(id uint32) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{Id: id})
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

func getEquipmentForCharacter(characterId uint32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId})
	}
}

func getEquipmentForCharacterBySlot(characterId uint32, slot int16) database.EntityProvider[entity] {
	return func(db *gorm.DB) model.Provider[entity] {
		return database.Query[entity](db, &entity{CharacterId: characterId, Slot: slot})
	}
}
