package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/item"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	DefaultInventoryCapacity uint32 = 24
)

func GetInventoryByType(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType string) (*Model, error) {
	return func(characterId uint32, inventoryType string) (*Model, error) {
		if it, ok := GetByteFromName(inventoryType); ok {
			return GetInventoryByTypeVal(l, db)(characterId, it)
		}
		return nil, errors.New("invalid inventory type")
	}
}

func GetInventoryByTypeFilterSlot(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType string, slot int16) (*Model, error) {
	return func(characterId uint32, inventoryType string, slot int16) (*Model, error) {
		if it, ok := GetByteFromName(inventoryType); ok {
			return GetInventoryByTypeValFilterSlot(l, db)(characterId, it, slot)
		}
		return nil, errors.New("invalid inventory type")
	}
}

func GetInventoryByTypeVal(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) (*Model, error) {
	return func(characterId uint32, inventoryType int8) (*Model, error) {
		i, err := get(db, characterId, inventoryType)
		if err != nil {
			return nil, err
		}

		if inventoryType == TypeValueEquip {
			i.items = getEquipInventoryItems(l, db)(characterId)
		} else {
			i.items = getInventoryItems(l, db)(characterId, inventoryType)
		}
		return i, nil
	}
}

func GetInventoryByTypeValFilterSlot(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, slot int16) (*Model, error) {
	return func(characterId uint32, inventoryType int8, slot int16) (*Model, error) {
		i, err := get(db, characterId, inventoryType)
		if err != nil {
			return nil, err
		}

		var items []Item
		if inventoryType == TypeValueEquip {
			items = getEquipInventoryItems(l, db)(characterId)
		} else {
			items = getInventoryItems(l, db)(characterId, inventoryType)
		}

		for _, it := range items {
			if it.Slot() == slot {
				i.items = append(i.items, it)
			}
		}

		return i, nil
	}
}

func getInventoryItems(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) []Item {
	return func(characterId uint32, inventoryType int8) []Item {
		results, err := item.GetForCharacterByInventory(l, db)(characterId, inventoryType)
		if err != nil {
			return make([]Item, 0)
		} else {
			var items = make([]Item, 0)
			for _, i := range results {
				items = append(items, Item{
					id:       i.Id(),
					itemType: ItemTypeItem,
					slot:     i.Slot(),
				})
			}
			return items
		}
	}
}

func getEquipInventoryItems(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) []Item {
	return func(characterId uint32) []Item {
		results, err := equipment.GetEquipmentForCharacter(l, db)(characterId)
		if err != nil {
			return make([]Item, 0)
		} else {
			var equips = make([]Item, 0)
			for _, e := range results {
				equips = append(equips, Item{
					id:       e.Id(),
					itemType: ItemTypeEquip,
					slot:     e.Slot(),
				})
			}
			return equips
		}
	}
}

func CreateInventory(db *gorm.DB) func(characterId uint32, inventoryType int8, capacity uint32) (*Model, error) {
	return func(characterId uint32, inventoryType int8, capacity uint32) (*Model, error) {
		return create(db, characterId, inventoryType, capacity)
	}
}

func CreateInitialInventories(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) error {
	return func(characterId uint32) error {
		_, err := CreateInventory(db)(characterId, TypeValueEquip, DefaultInventoryCapacity)
		if err != nil {
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueUse, DefaultInventoryCapacity)
		if err != nil {
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueSetup, DefaultInventoryCapacity)
		if err != nil {
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueETC, DefaultInventoryCapacity)
		if err != nil {
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueCash, DefaultInventoryCapacity)
		if err != nil {
			return err
		}
		return nil
	}
}
