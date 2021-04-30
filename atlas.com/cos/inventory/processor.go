package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/item"
	"errors"
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

const (
	DefaultInventoryCapacity uint32 = 24
)

func (p processor) GetInventoryByType(characterId uint32, inventoryType string) (*Model, error) {
	if it, ok := GetByteFromName(inventoryType); ok {
		return p.GetInventoryByTypeVal(characterId, it)
	}
	return nil, errors.New("invalid inventory type")
}

func (p processor) GetInventoryByTypeFilterSlot(characterId uint32, inventoryType string, slot int16) (*Model, error) {
	if it, ok := GetByteFromName(inventoryType); ok {
		return p.GetInventoryByTypeValFilterSlot(characterId, it, slot)
	}
	return nil, errors.New("invalid inventory type")
}

func (p processor) GetInventoryByTypeVal(characterId uint32, inventoryType int8) (*Model, error) {
	i, err := Get(p.db, characterId, inventoryType)
	if err != nil {
		return nil, err
	}

	if inventoryType == TypeValueEquip {
		i.items = p.getEquipInventoryItems(characterId)
	} else {
		i.items = p.getInventoryItems(characterId, inventoryType)
	}
	return i, nil
}

func (p processor) GetInventoryByTypeValFilterSlot(characterId uint32, inventoryType int8, slot int16) (*Model, error) {
	i, err := Get(p.db, characterId, inventoryType)
	if err != nil {
		return nil, err
	}

	var items []InventoryItem
	if inventoryType == TypeValueEquip {
		items = p.getEquipInventoryItems(characterId)
	} else {
		items = p.getInventoryItems(characterId, inventoryType)
	}

	for _, it := range items {
		if it.Slot() == slot {
			i.items = append(i.items, it)
		}
	}

	return i, nil
}

func (p processor) getInventoryItems(characterId uint32, inventoryType int8) []InventoryItem {
	results, err := item.Processor(p.l, p.db).GetForCharacterByInventory(characterId, inventoryType)
	if err != nil {
		return make([]InventoryItem, 0)
	} else {
		var items = make([]InventoryItem, 0)
		for _, i := range results {
			items = append(items, InventoryItem{
				id:       i.Id(),
				itemType: ItemTypeItem,
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

func (p processor) CreateInventory(characterId uint32, inventoryType int8, capacity uint32) (*Model, error) {
	return Create(p.db, characterId, inventoryType, capacity)
}

func (p processor) CreateInitialInventories(characterId uint32) error {
	_, err := p.CreateInventory(characterId, TypeValueEquip, DefaultInventoryCapacity)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueUse, DefaultInventoryCapacity)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueSetup, DefaultInventoryCapacity)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueETC, DefaultInventoryCapacity)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueCash, DefaultInventoryCapacity)
	if err != nil {
		return err
	}
	return nil
}
