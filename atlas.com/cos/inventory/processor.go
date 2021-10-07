package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	DefaultInventoryCapacity uint32 = 24
)

func GetInventory(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType string, filters ...ItemFilter) (*Model, error) {
	return func(characterId uint32, inventoryType string, filters ...ItemFilter) (*Model, error) {
		if it, ok := GetByteFromName(inventoryType); ok {
			return GetInventoryByTypeVal(l, db)(characterId, it, filters...)
		}
		return nil, errors.New("invalid inventory type")
	}
}

func GetInventoryByTypeVal(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, filters ...ItemFilter) (*Model, error) {
	return func(characterId uint32, inventoryType int8, filters ...ItemFilter) (*Model, error) {
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
			ok := true
			for _, filter := range filters {
				if !filter(it) {
					ok = false
					break
				}
			}
			if ok {
				i.items = append(i.items, it)
			}
		}

		return i, nil
	}
}

type ItemFilter func(i Item) bool

func FilterSlot(slot int16) ItemFilter {
	return func(i Item) bool {
		return i.Slot() == slot
	}
}

func FilterItemId(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(itemId uint32) ItemFilter {
	return func(itemId uint32) ItemFilter {
		return func(i Item) bool {
			if i.ItemType() == ItemTypeItem {
				ii, err := item.GetItemById(l, db)(i.Id())
				if err != nil {
					return false
				}
				return ii.ItemId() == itemId
			} else {
				ee, err := equipment.GetEquipmentById(l, db)(i.Id())
				if err != nil {
					return false
				}
				es, err := statistics.GetEquipmentStatistics(l, span)(ee.EquipmentId())
				if err != nil {
					return false
				}
				return es.ItemId() == itemId
			}
		}
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
		l.Debugf("Creating default inventories for character %d.", characterId)
		_, err := CreateInventory(db)(characterId, TypeValueEquip, DefaultInventoryCapacity)
		if err != nil {
			l.WithError(err).Errorf("Error creating default EQUIP inventory for character %d.", characterId)
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueUse, DefaultInventoryCapacity)
		if err != nil {
			l.WithError(err).Errorf("Error creating default USE inventory for character %d.", characterId)
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueSetup, DefaultInventoryCapacity)
		if err != nil {
			l.WithError(err).Errorf("Error creating default SETUP inventory for character %d.", characterId)
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueETC, DefaultInventoryCapacity)
		if err != nil {
			l.WithError(err).Errorf("Error creating default ETC inventory for character %d.", characterId)
			return err
		}
		_, err = CreateInventory(db)(characterId, TypeValueCash, DefaultInventoryCapacity)
		if err != nil {
			l.WithError(err).Errorf("Error creating default CASH inventory for character %d.", characterId)
			return err
		}
		return nil
	}
}
