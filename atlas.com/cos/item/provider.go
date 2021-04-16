package item

import (
	"gorm.io/gorm"
	"sort"
)

func GetItemsForCharacter(db *gorm.DB, characterId uint32, inventoryType byte, itemId uint32) ([]*Model, error) {
	var results []entity
	err := db.First(&results).Where("characterId = ? AND inventoryType = ? AND itemId = ?", characterId, inventoryType, itemId).Error
	if err != nil {
		return nil, err
	}

	var items = make([]*Model, 0)
	for _, e := range results {
		items = append(items, makeItem(&e))
	}
	return items, nil
}

func GetItemsForCharacterByInventory(db *gorm.DB, characterId uint32, inventoryType byte) ([]*Model, error) {
	var results []entity
	err := db.First(&results).Where("characterId = ? AND inventoryType = ?", characterId, inventoryType).Error
	if err != nil {
		return nil, err
	}

	var items = make([]*Model, 0)
	for _, e := range results {
		items = append(items, makeItem(&e))
	}
	return items, nil
}

func GetNextFreeEquipmentSlot(db *gorm.DB, characterId uint32, inventoryType byte) (int16, error) {
	items, err := GetItemsForCharacterByInventory(db, characterId, inventoryType)
	if err != nil {
		return 0, err
	}
	if len(items) == 0 {
		return 0, nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Slot() > items[i].Slot()
	})
	return minFreeSlot(items), nil
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

func makeItem(e *entity) *Model {
	return &Model{
		id:            e.id,
		characterId:   e.characterId,
		inventoryType: e.inventoryType,
		itemId:        e.itemId,
		quantity:      e.quantity,
		slot:          e.slot,
	}
}
