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

func (p processor) GetInventoryByType(characterId uint32, inventoryType string) (*Model, error) {
	if it, ok := GetByteFromName(inventoryType); ok {
		return p.GetInventoryByTypeVal(characterId, it)
	}
	return nil, errors.New("invalid inventory type")
}

func (p processor) GetInventoryByTypeVal(characterId uint32, inventoryType byte) (*Model, error) {
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

func (p processor) getInventoryItems(characterId uint32, inventoryType byte) []InventoryItem {
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

func (p processor) CreateInventory(characterId uint32, inventoryType byte, capacity uint32) (*Model, error) {
	return Create(p.db, characterId, inventoryType, capacity)
}

func (p processor) CreateInitialInventories(characterId uint32) error {
	_, err := p.CreateInventory(characterId, TypeValueEquip, 4)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueUse, 4)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueSetup, 4)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueETC, 4)
	if err != nil {
		return err
	}
	_, err = p.CreateInventory(characterId, TypeValueCash, 4)
	if err != nil {
		return err
	}
	return nil
}
