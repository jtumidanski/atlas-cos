package item

import (
	"gorm.io/gorm"
	"sort"
)

func GetById(db *gorm.DB, id uint32) (*Model, error) {
	var result entity
	err := db.Where(&entity{Id: id}).First(&result).Error
	if err != nil {
		return nil, err
	}

	return makeItem(&result), nil
}

func GetItemsForCharacter(db *gorm.DB, characterId uint32, inventoryType byte, itemId uint32) ([]*Model, error) {
	var results []entity
	err := db.Where(&entity{CharacterId: characterId, InventoryType: inventoryType, ItemId: itemId}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var items = make([]*Model, 0)
	for _, e := range results {
		items = append(items, makeItem(&e))
	}
	return items, nil
}

func GetForCharacterByInventory(db *gorm.DB, characterId uint32, inventoryType byte) ([]*Model, error) {
	var results []entity
	err := db.Where(&entity{CharacterId: characterId, InventoryType: inventoryType}).Find(&results).Error
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
	err := db.Where(&entity{CharacterId: characterId, InventoryType: inventoryType}).Find(&results).Error
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
		id:            e.Id,
		characterId:   e.CharacterId,
		inventoryType: e.InventoryType,
		itemId:        e.ItemId,
		quantity:      e.Quantity,
		slot:          e.Slot,
	}
}
