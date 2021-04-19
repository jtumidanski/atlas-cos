package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/item"
	"errors"
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetInventoryByType(characterId uint32, inventoryType string) (*Model, error) {
	if it, ok := GetByteFromName(inventoryType); ok {
		return p.getInventoryByType(characterId, it)
	}
	return nil, errors.New("invalid inventory type")
}

func (p processor) getInventoryByType(characterId uint32, inventoryType byte) (*Model, error) {
	var items = make([]InventoryItem, 0)
	if inventoryType == TypeValueEquip {
		items = p.getEquipInventoryItems(characterId)
	} else {
		items = p.getInventoryItems(characterId, inventoryType)
	}

	name, _ := GetTypeFromByte(inventoryType)
	return &Model{
		id:            inventoryType,
		inventoryType: name,
		capacity:      4,
		items:         items,
	}, nil
}

func (p processor) getInventoryItems(characterId uint32, inventoryType byte) []InventoryItem {
	results, err := item.Processor(p.l, p.db).GetForCharacterByInventory(characterId, inventoryType)
	if err != nil {
		return make([]InventoryItem, 0)
	} else {
		var items = make([]InventoryItem, 0)
		for _, i := range results {
			items = append(items, InventoryItem{
				id:       i.Id(),
				itemType: ItemTypeEquip,
				slot:     i.Slot(),
			})
		}
		return items
	}
}

func (p processor) getEquipInventoryItems(characterId uint32) []InventoryItem {
	results, err := equipment.Processor(p.l, p.db).GetEquipmentForCharacter(characterId)
	if err != nil {
		return make([]InventoryItem, 0)
	} else {
		var equips = make([]InventoryItem, 0)
		for _, e := range results {
			equips = append(equips, InventoryItem{
				id:       e.Id(),
				itemType: ItemTypeEquip,
				slot:     e.Slot(),
			})
		}
		return equips
	}
}
