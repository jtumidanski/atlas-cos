package item

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetItemsForCharacter(characterId uint32, inventoryType int8, itemId uint32) []*Model {
	items, err := GetItemsForCharacter(p.db, characterId, inventoryType, itemId)
	if err != nil {
		return make([]*Model, 0)
	}
	return items
}

func (p processor) GetForCharacterByInventory(characterId uint32, inventoryType int8) ([]*Model, error) {
	return GetForCharacterByInventory(p.db, characterId, inventoryType)
}

func (p processor) UpdateItemQuantity(id uint32, quantity uint32) error {
	return Update(p.db, id, SetQuantity(quantity))
}

func (p processor) CreateItemForCharacter(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (*Model, error) {
	slot, err := GetNextFreeEquipmentSlot(p.db, characterId, inventoryType)
	if err != nil {
		return nil, err
	}
	return CreateItemForCharacter(p.db, characterId, inventoryType, itemId, quantity, slot)
}

func (p processor) GetItemById(id uint32) (*Model, error) {
	return GetById(p.db, id)
}
