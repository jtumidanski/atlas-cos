package item

import (
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

func (p processor) GetItemsForCharacter(characterId uint32, inventoryType byte, itemId uint32) []*Model {
	items, err := GetItemsForCharacter(p.db, characterId, inventoryType, itemId)
	if err != nil {
		return make([]*Model, 0)
	}
	return items
}

func (p processor) UpdateItemQuantity(id uint32, quantity uint32) error {
	return Update(p.db, id, SetQuantity(quantity))
}

func (p processor) CreateItemForCharacter(characterId uint32, inventoryType byte, itemId uint32, quantity uint32) (*Model, error) {
	slot, err := GetNextFreeEquipmentSlot(p.db, characterId, inventoryType)
	if err != nil {
		return nil, err
	}
	return CreateItemForCharacter(p.db, characterId, inventoryType, itemId, quantity, slot)
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}
