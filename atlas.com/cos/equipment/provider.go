package equipment

import (
	"gorm.io/gorm"
	"sort"
)

func getEquipment(db *gorm.DB, query interface{}) (*Model, error) {
	var result entity
	err := db.Where(query).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return makeEquipment(&result), nil
}

func getEquipments(db *gorm.DB, query interface{}) ([]*Model, error) {
	var results []entity
	err := db.Where(query).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var equipment = make([]*Model, 0)
	for _, e := range results {
		equipment = append(equipment, makeEquipment(&e))
	}
	return equipment, nil
}

func getByEquipmentId(db *gorm.DB, equipmentId uint32) (*Model, error) {
	return getEquipment(db, &entity{EquipmentId: equipmentId})
}

func getById(db *gorm.DB, id uint32) (*Model, error) {
	return getEquipment(db, &entity{Id: id})
}

func getNextFreeEquipmentSlot(db *gorm.DB, characterId uint32) (int16, error) {
	equipment, err := getEquipmentForCharacter(db, characterId)
	if err != nil {
		return 1, err
	}
	if len(equipment) == 0 {
		return 1, nil
	}

	sort.Slice(equipment, func(i, j int) bool {
		return equipment[i].Slot() < equipment[j].Slot()
	})
	return minFreeSlot(equipment), nil
}

func minFreeSlot(items []*Model) int16 {
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
		} else if items[i].Slot() < 0 {
			i += 1
		}
	}
}

func getEquipmentForCharacter(db *gorm.DB, characterId uint32) ([]*Model, error) {
	return getEquipments(db, &entity{CharacterId: characterId})
}

func getEquipmentForCharacterBySlot(db *gorm.DB, characterId uint32, slot int16) (*Model, error) {
	return getEquipment(db, &entity{CharacterId: characterId, Slot: slot})
}
