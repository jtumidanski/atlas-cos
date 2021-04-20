package equipment

import (
	"gorm.io/gorm"
	"sort"
)

func GetByEquipmentId(db *gorm.DB, equipmentId uint32) (*Model, error) {
	var result entity
	err := db.Where(&entity{EquipmentId: equipmentId}).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return makeEquipment(&result), nil
}

func GetById(db *gorm.DB, id uint32) (*Model, error) {
	var result entity
	err := db.Where(&entity{Id: id}).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return makeEquipment(&result), nil
}

func GetNextFreeEquipmentSlot(db *gorm.DB, characterId uint32) (int16, error) {
	equipment, err := GetEquipmentForCharacter(db, characterId)
	if err != nil {
		return 0, err
	}
	if len(equipment) == 0 {
		return 0, nil
	}

	sort.Slice(equipment, func(i, j int) bool {
		return equipment[i].Slot() > equipment[i].Slot()
	})
	return minFreeSlot(equipment), nil
}

func minFreeSlot(items []*Model) int16 {
	slot := int16(0)
	i := 0

	for {
		if i >= len(items) {
			return slot
		} else if slot < items[i].Slot() {
			return slot
		} else if slot == items[i].Slot() {
			slot += 1
			i += 1
		} else if items[i].Slot() < 0 {
			i += 1
		}
	}
}

func GetEquipmentForCharacter(db *gorm.DB, characterId uint32) ([]*Model, error) {
	var results []entity
	err := db.Where(&entity{CharacterId: characterId}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var equipment = make([]*Model, 0)
	for _, e := range results {
		equipment = append(equipment, makeEquipment(&e))
	}
	return equipment, nil
}

func GetEquipmentForCharacterBySlot(db *gorm.DB, characterId uint32, slot int16) (*Model, error) {
	var results entity
	err := db.Where(&entity{CharacterId: characterId, Slot: slot}, characterId, slot).First(&results).Error
	if err != nil {
		return nil, err
	}
	return makeEquipment(&results), nil
}