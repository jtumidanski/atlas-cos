package item

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sort"
)

func getById(db *gorm.DB, id uint32) (*Model, error) {
	var result entity
	err := db.Where(&entity{Id: id}).First(&result).Error
	if err != nil {
		return nil, err
	}

	return makeItem(&result), nil
}

func getItemsForCharacter(db *gorm.DB, characterId uint32, inventoryType int8, itemId uint32) ([]*Model, error) {
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

func getItemForCharacter(db *gorm.DB, characterId uint32, inventoryType int8, slot int16) (*Model, error) {
	var results entity
	err := db.Where(&entity{CharacterId: characterId, InventoryType: inventoryType, Slot: slot}).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return makeItem(&results), nil
}

func getForCharacterByInventory(db *gorm.DB, characterId uint32, inventoryType int8) ([]*Model, error) {
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

func getItemsForCharacterByInventory(db *gorm.DB, characterId uint32, inventoryType int8) ([]*Model, error) {
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

func getNextFreeEquipmentSlot(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) (int16, error) {
	return func(characterId uint32, inventoryType int8) (int16, error) {
		items, err := getItemsForCharacterByInventory(db, characterId, inventoryType)
		if err != nil {
			return 1, err
		}

		l.Debugf("Character %d has %d items in %d inventory.", characterId, len(items), inventoryType)

		if len(items) == 0 {
			l.Debugf("Defaulting item slot choice to 1 for character %d.", characterId)
			return 1, nil
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Slot() < items[j].Slot()
		})
		slot := minFreeSlot(items)

		l.Debugf("Chose slot %d for new item in inventory %d for character %d.", slot, inventoryType, characterId)

		return slot, nil
	}
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
		} else if items[i].Slot() <= 0 {
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
