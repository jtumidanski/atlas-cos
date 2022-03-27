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

func DropFromInventorySlot(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, inventoryType int8, slot int16, quantity int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, inventoryType int8, slot int16, quantity int16) error {
		if slot < 0 {
			return dropEquippedItem(l, db, span)(worldId, channelId, mapId, characterId, x, y, slot)
		}

		if inventoryType == TypeValueEquip {
			return dropEquipItem(l, db, span)(worldId, channelId, mapId, characterId, x, y, slot)
		}

		return dropItem(l, db, span)(worldId, channelId, mapId, characterId, x, y, inventoryType, slot, quantity)
	}
}

func dropItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, inventoryType int8, slot int16, quantity int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, inventoryType int8, slot int16, quantity int16) error {
		adjustment, err := item.DropItem(l, db, span)(worldId, channelId, characterId, inventoryType, slot, quantity)
		if err != nil {
			l.WithError(err).Errorf("Unable to drop item from inventory %d slot %d for character %d.", inventoryType, slot, characterId)
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
		emitSpawnItemDropCommand(l, span)(worldId, channelId, mapId, adjustment.ItemId(), uint32(quantity), 0, 0, x, y, characterId, 0)
		return nil
	}
}

func dropEquipItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
		eid, err := equipment.DropEquipment(l, db, span)(worldId, channelId, characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to drop equip item from slot %d for character %d.", slot, characterId)
			return err
		}

		itemId, err := equipment.GetItemIdForEquipment(l, span)(eid)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment %d.", eid)
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, 3, itemId, 1, 1, slot, 0)

		emitSpawnEquipDropCommand(l, span)(worldId, channelId, mapId, itemId, eid, 0, x, y, characterId, 0)
		return nil
	}
}

func dropEquippedItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
		eid, err := equipment.DropEquippedItem(l, db, span)(worldId, channelId, characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to drop equipped item from slot %d for character %d.", slot, characterId)
			return err
		}

		itemId, err := equipment.GetItemIdForEquipment(l, span)(eid)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment %d.", eid)
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, 3, itemId, 1, 1, slot, 0)

		emitSpawnEquipDropCommand(l, span)(worldId, channelId, mapId, itemId, eid, 0, x, y, characterId, 0)
		return nil
	}
}

func Move(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) error {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) error {
		if inventoryType == TypeValueEquip {
			return moveEquipItem(l, db, span)(characterId, inventoryType, source, destination)
		}
		return moveItem(l, db, span)(characterId, inventoryType, source, destination)
	}
}

func moveItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) error {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) error {
		adjustment, err := item.MoveItem(l, db, span)(characterId, inventoryType, source, destination)
		if err != nil {
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
		return nil
	}
}

func moveEquipItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) error {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) error {
		adjustment, err := equipment.MoveItem(l, db, span)(characterId, source, destination)
		if err != nil {
			return err
		}
		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
		return nil
	}
}

func EquipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32) {
	return func(characterId uint32, equipmentId uint32) {
		adjustment, err := equipment.EquipItemForCharacter(l, db, span)(characterId, equipmentId)
		if err != nil {
			return
		}
		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
	}
}

func CreateAndEquip(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, items ...uint32) {
	return func(characterId uint32, items ...uint32) {
		for _, i := range items {
			createAndEquip(l, db, span)(characterId, i)
		}
	}
}

func createAndEquip(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, itemId uint32) {
	return func(characterId uint32, itemId uint32) {
		eid, err := statistics.Create(l, span)(itemId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
			return
		}

		if e, err := equipment.CreateForCharacter(l, db)(characterId, itemId, eid, true); err == nil {
			EquipItemForCharacter(l, db, span)(characterId, e.EquipmentId())
		} else {
			l.Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
		}
	}
}

func UnequipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32, oldSlot int16) {
	return func(characterId uint32, equipmentId uint32, oldSlot int16) {
		adjustment, err := equipment.UnequipItemForCharacter(l, db, span)(characterId, equipmentId, oldSlot)
		if err != nil {
			return
		}
		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
	}
}

func GainEquipment(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, itemId uint32, equipmentId uint32) error {
	return func(characterId uint32, itemId uint32, equipmentId uint32) error {
		adjustment, err := equipment.GainItem(l, db, span)(characterId, itemId, equipmentId)
		if err != nil {
			return err
		}
		emitInventoryModificationEvent(l, span)(characterId, true, adjustment.Mode(), adjustment.ItemId(), adjustment.InventoryType(), adjustment.Quantity(), adjustment.Slot(), adjustment.OldSlot())
		return nil
	}
}

func GainItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) error {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) error {
		events, err := item.GainItem(l, db, span)(characterId, inventoryType, itemId, quantity)
		if err != nil {
			return err
		}
		for _, e := range events {
			emitInventoryModificationEvent(l, span)(characterId, true, e.Mode(), e.ItemId(), e.InventoryType(), e.Quantity(), e.Slot(), e.OldSlot())
		}
		return nil
	}
}

func LoseItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, itemId uint32, quantity int32) error {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity int32) error {
		events, err := item.LoseItem(l, db, span)(characterId, inventoryType, itemId, quantity)
		if err != nil {
			return err
		}
		for _, e := range events {
			emitInventoryModificationEvent(l, span)(characterId, true, e.Mode(), e.ItemId(), e.InventoryType(), e.Quantity(), e.Slot(), e.OldSlot())
		}
		return nil
	}
}
