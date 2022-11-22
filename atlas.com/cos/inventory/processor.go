package inventory

import (
	"atlas-cos/database"
	"atlas-cos/equipment/statistics"
	"atlas-cos/inventory/item"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

const (
	DefaultInventoryCapacity uint32 = 24
)

type adjustment struct {
	mode          byte
	itemId        uint32
	inventoryType int8
	quantity      uint32
	slot          int16
	oldSlot       int16
}

func (i adjustment) Mode() byte {
	return i.mode
}

func (i adjustment) ItemId() uint32 {
	return i.itemId
}

func (i adjustment) InventoryType() int8 {
	return i.inventoryType
}

func (i adjustment) Quantity() uint32 {
	return i.quantity
}

func (i adjustment) Slot() int16 {
	return i.slot
}

func (i adjustment) OldSlot() int16 {
	return i.oldSlot
}

func GetInventory(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType string, filters ...ItemFilter) (Model, error) {
	return func(characterId uint32, inventoryType string, filters ...ItemFilter) (Model, error) {
		if it, ok := GetByteFromName(inventoryType); ok {
			return GetInventoryByTypeVal(l, db)(characterId, it, filters...)
		}
		return Model{}, errors.New("invalid inventory type")
	}
}

func GetInventoryByTypeVal(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, filters ...ItemFilter) (Model, error) {
	return func(characterId uint32, inventoryType int8, filters ...ItemFilter) (Model, error) {
		i, err := database.ModelProvider[Model, entity](db)(get(characterId, inventoryType), makeInventory)()
		if err != nil {
			return Model{}, err
		}

		items, err := item.GetByInventory(l, db)(i.Id())
		if err != nil {
			return Model{}, err
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

func GetEquipment(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]item.EquipmentModel, error) {
	return func(characterId uint32) ([]item.EquipmentModel, error) {
		i, err := GetInventoryByTypeVal(l, db)(characterId, TypeValueEquip)
		if err != nil {
			return nil, err
		}
		return item.GetEquipment(l, db)(i.Id())
	}
}

func GetEquippedItemBySlot(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, slot int16) (item.EquipmentModel, error) {
	return func(characterId uint32, slot int16) (item.EquipmentModel, error) {
		i, err := GetInventoryByTypeVal(l, db)(characterId, TypeValueEquip)
		if err != nil {
			return item.EquipmentModel{}, err
		}
		return item.GetEquippedItemBySlot(l, db)(i.Id(), slot)
	}
}

type ItemFilter func(i item.Model) bool

func FilterSlot(slot int16) ItemFilter {
	return func(i item.Model) bool {
		return i.Slot() == slot
	}
}

func FilterItemId(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(itemId uint32) ItemFilter {
	return func(itemId uint32) ItemFilter {
		return func(i item.Model) bool {
			ii, err := item.GetById(l, db)(i.Id())
			if err != nil {
				return false
			}
			return ii.ItemId() == itemId
		}
	}
}

func CreateInventory(db *gorm.DB) func(characterId uint32, inventoryType int8, capacity uint32) (Model, error) {
	return func(characterId uint32, inventoryType int8, capacity uint32) (Model, error) {
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
		l.Debugf("Character %d dropping %d item in inventory %d slot %d.", characterId, quantity, inventoryType, slot)
		ii, err := GetInventoryByTypeVal(l, db)(characterId, inventoryType)
		if err != nil {
			return err
		}

		i, err := item.GetBySlot(l, db)(ii.Id(), slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return err
		}

		initialQuantity := uint32(1)
		if val, ok := i.(*item.ItemModel); ok {
			initialQuantity = val.Quantity()
		}

		if initialQuantity <= uint32(quantity) {
			err = item.RemoveItem(l, db)(characterId, i.Id())
			if err != nil {
				l.WithError(err).Errorf("Could not remove item %d from character %d inventory.", i.Id(), characterId)
				return err
			}
			emitInventoryModificationEvent(l, span)(characterId, true, 3, i.ItemId(), inventoryType, uint32(quantity), slot, 0)
			emitSpawnItemDropCommand(l, span)(worldId, channelId, mapId, i.ItemId(), uint32(quantity), 0, 0, x, y, characterId, 0)
		} else {
			newQuantity := initialQuantity - uint32(quantity)
			err = item.UpdateQuantity(l, db)(i.Id(), newQuantity)
			if err != nil {
				l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
				return err
			}
			emitInventoryModificationEvent(l, span)(characterId, true, 1, i.ItemId(), inventoryType, newQuantity, slot, 0)
			emitSpawnItemDropCommand(l, span)(worldId, channelId, mapId, i.ItemId(), uint32(quantity), 0, 0, x, y, characterId, 0)
		}
		return nil
	}
}

func dropEquipItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
		l.Debugf("Character %d dropping equipment in slot %d.", characterId, slot)
		e, err := GetEquippedItemBySlot(l, db)(characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return err
		}

		err = item.RemoveItem(l, db)(characterId, e.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to remove item %d.", e.Id())
			return err
		}

		emitItemUnequipped(l, span)(characterId)
		emitInventoryModificationEvent(l, span)(characterId, true, 3, e.ItemId(), 1, 1, slot, 0)
		emitSpawnEquipDropCommand(l, span)(worldId, channelId, mapId, e.ItemId(), e.EquipmentId(), 0, x, y, characterId, 0)
		return nil
	}
}

func dropEquippedItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, x int16, y int16, slot int16) error {
		l.Debugf("Character %d dropping equipment in slot %d.", characterId, slot)
		e, err := GetEquippedItemBySlot(l, db)(characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return err
		}

		err = item.RemoveItem(l, db)(characterId, e.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to remove item %d.", e.Id())
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, 3, e.ItemId(), 1, 1, slot, 0)
		emitSpawnEquipDropCommand(l, span)(worldId, channelId, mapId, e.ItemId(), e.EquipmentId(), 0, x, y, characterId, 0)
		return nil
	}
}

func Move(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) error {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) error {
		itemId := uint32(0)
		quantity := uint32(1)
		txError := db.Transaction(func(tx *gorm.DB) error {
			ii, err := GetInventoryByTypeVal(l, tx)(characterId, inventoryType)
			if err != nil {
				return err
			}

			i, err := item.GetBySlot(l, tx)(ii.Id(), source)
			if err != nil || i.Id() == 0 {
				l.Warnf("Item movement requested, but no equipment for character %d in slot %d.", characterId, source)
				return err
			}
			itemId = i.ItemId()
			if val, ok := i.(*item.ItemModel); ok {
				quantity = val.Quantity()
			}

			temporarySlot := int16(math.MinInt16)
			otherItem, err := item.GetBySlot(l, tx)(ii.Id(), destination)
			if err == nil && otherItem.Id() != 0 {
				l.Debugf("Item %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", otherItem.Id(), destination, temporarySlot, characterId)
				err = item.UpdateSlot(l, tx)(otherItem.Id(), temporarySlot)
				if err != nil {
					l.WithError(err).Errorf("Unable to move other item %d to temporary slot %d.", otherItem.Id(), temporarySlot)
					return err
				}
			}

			err = item.UpdateSlot(l, tx)(i.Id(), destination)
			if err != nil {
				return err
			}
			l.Debugf("Moved item %d from slot %d to %d for character %d.", i.Id(), source, destination, characterId)

			if otherItem != nil && otherItem.Id() != 0 {
				err = item.UpdateSlot(l, tx)(otherItem.Id(), source)
				if err != nil {
					l.WithError(err).Errorf("Unable to move other item %d to resulting slot %d.", otherItem.Id(), source)
					return err
				}
				l.Debugf("Moved item %d from slot %d to %d for character %d.", otherItem.Id(), temporarySlot, source, characterId)
			}
			return nil
		})
		if txError != nil {
			return txError
		}

		emitInventoryModificationEvent(l, span)(characterId, true, 2, itemId, inventoryType, quantity, destination, source)
		return nil
	}
}

func EquipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32) {
	return func(characterId uint32, equipmentId uint32) {
		l.Debugf("Received request to equip %d for character %d.", equipmentId, characterId)
		e, err := item.GetByEquipmentId(l, db)(equipmentId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
			return
		}

		l.Debugf("Equipment %d is item %d for character %d.", equipmentId, e.ItemId(), characterId)

		slots, err := item.GetEquipmentSlotDestination(l, span)(e.ItemId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve destination slots for item %d.", e.ItemId())
			return
		} else if len(slots) <= 0 {
			l.Errorf("Unable to retrieve destination slots for item %d. %s.", e.ItemId())
			return
		}
		slot := slots[0]
		l.Debugf("Equipment %d to be equipped in slot %d for character %d.", equipmentId, slot, characterId)

		temporarySlot := int16(math.MinInt16)

		existingSlot := e.Slot()
		err = db.Transaction(func(tx *gorm.DB) error {
			if equip, err := GetEquippedItemBySlot(l, tx)(characterId, slot); err == nil && equip.EquipmentId() != 0 {
				l.Debugf("Equipment %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", equip.EquipmentId(), slot, temporarySlot, characterId)
				_ = item.UpdateSlot(l, tx)(equip.Id(), temporarySlot)
			}

			err = item.UpdateSlot(l, tx)(e.Id(), slot)
			if err != nil {
				return err
			}
			l.Debugf("Moved item %d from slot %d to %d for character %d.", e.ItemId(), existingSlot, slot, characterId)

			if equip, err := GetEquippedItemBySlot(l, tx)(characterId, temporarySlot); err == nil && equip.EquipmentId() != 0 {
				err := item.UpdateSlot(l, tx)(equip.Id(), existingSlot)
				if err != nil {
					return err
				}
				l.Debugf("Moved item from temporary location %d to slot %d for character %d.", temporarySlot, existingSlot, characterId)
			}
			return nil
		})
		if err != nil {
			l.WithError(err).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
			return
		}

		emitItemEquipped(l, span)(characterId)
		emitInventoryModificationEvent(l, span)(characterId, true, 2, e.ItemId(), TypeValueEquip, 1, slot, existingSlot)
	}
}

func CreateItem(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (item.Model, error) {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (item.Model, error) {
		inv, err := GetInventoryByTypeVal(l, db)(characterId, inventoryType)
		if err != nil {
			return item.ItemModel{}, err
		}
		return item.CreateItem(l, db)(characterId, inv.Id(), inventoryType, itemId, quantity)
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

		i, err := GetInventoryByTypeVal(l, db)(characterId, TypeValueEquip)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve inventory for character %d.", characterId)
			return
		}

		if e, err := item.CreateEquipment(l, db)(characterId, i.Id(), itemId, eid, true); err == nil {
			EquipItemForCharacter(l, db, span)(characterId, e.EquipmentId())
		} else {
			l.Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
		}
	}
}

func UnequipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32, oldSlot int16) {
	return func(characterId uint32, equipmentId uint32, oldSlot int16) {
		l.Debugf("Received request to unequip %d for character %d.", equipmentId, characterId)

		itemId := uint32(0)
		newSlot := int16(0)
		txErr := db.Transaction(func(tx *gorm.DB) error {
			e, err := item.GetByEquipmentId(l, tx)(equipmentId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
				return err
			}

			ii, err := GetInventoryByTypeVal(l, tx)(characterId, TypeValueEquip)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate inventory for character %d.", characterId)
				return err
			}

			newSlot, err = item.GetNextFreeSlot(l, tx)(ii.Id())
			if err != nil {
				l.WithError(err).Errorf("Unable to get next free equipment slot")
				return err
			}

			err = item.UpdateSlot(l, tx)(e.Id(), newSlot)
			if err != nil {
				return err
			}

			l.Debugf("Unequipped %d for character %d and place it in slot %d, from %d.", equipmentId, characterId, newSlot, oldSlot)
			return nil
		})
		if txErr != nil {
			l.WithError(txErr).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
			return
		}
		emitItemUnequipped(l, span)(characterId)
		emitInventoryModificationEvent(l, span)(characterId, true, 2, itemId, TypeValueEquip, 1, newSlot, oldSlot)
	}
}

// Deprecated: going to do a generic GainItem call instead
func GainEquipment(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, itemId uint32, equipmentId uint32) error {
	return func(characterId uint32, itemId uint32, equipmentId uint32) error {
		//TODO verify inventory space
		ii, err := GetInventoryByTypeVal(l, db)(characterId, TypeValueEquip)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate inventory for character %d.", characterId)
			return err
		}

		e, err := item.CreateEquipment(l, db)(characterId, ii.Id(), itemId, equipmentId, false)
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
			return err
		}

		emitInventoryModificationEvent(l, span)(characterId, true, 0, itemId, TypeValueEquip, 1, e.Slot(), 0)
		return nil
	}
}

func GainItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) error {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) error {
		//TODO verify inventory space
		slotMax := item.MaxInSlot()
		runningQuantity := quantity

		var events = make([]adjustment, 0)

		txError := db.Transaction(func(tx *gorm.DB) error {
			ii, err := GetInventoryByTypeVal(l, tx)(characterId, inventoryType)
			if err != nil {
				return err
			}

			existingItems, err := item.GetByItemId(l, tx)(ii.Id(), itemId)
			if err != nil {
				return err
			}
			// breaks for a rechargeable item.
			if len(existingItems) > 0 {
				index := 0
				for runningQuantity > 0 {
					if index < len(existingItems) {
						i := existingItems[index]
						oldQuantity := uint32(1)
						if val, ok := i.(*item.ItemModel); ok {
							oldQuantity = val.Quantity()
						}

						if oldQuantity < slotMax {
							newQuantity := uint32(math.Min(float64(oldQuantity+runningQuantity), float64(slotMax)))
							runningQuantity = runningQuantity - (newQuantity - oldQuantity)
							err := item.UpdateQuantity(l, tx)(i.Id(), newQuantity)
							if err != nil {
								l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
							} else {
								events = append(events, adjustment{mode: 1, itemId: itemId, inventoryType: inventoryType, quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
							}
						}
						index++
					} else {
						break
					}
				}
			}
			for runningQuantity > 0 {
				newQuantity := uint32(math.Min(float64(runningQuantity), float64(slotMax)))
				runningQuantity = runningQuantity - newQuantity
				i, err := item.CreateItem(l, tx)(characterId, ii.Id(), inventoryType, itemId, newQuantity)
				if err != nil {
					l.WithError(err).Errorf("Unable to create item %d that character %d picked up.", itemId, characterId)
					return err
				}
				events = append(events, adjustment{mode: 0, itemId: itemId, inventoryType: inventoryType, quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
			}
			return nil
		})
		if txError != nil {
			return txError
		}
		for _, e := range events {
			emitInventoryModificationEvent(l, span)(characterId, true, e.Mode(), e.ItemId(), e.InventoryType(), e.Quantity(), e.Slot(), e.OldSlot())
		}
		return nil
	}
}

func LoseItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inventoryType int8, itemId uint32, quantity int32) error {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity int32) error {
		//TODO verify inventory space

		runningQuantity := uint32(quantity * -1)

		var events = make([]adjustment, 0)

		txError := db.Transaction(func(tx *gorm.DB) error {
			ii, err := GetInventoryByTypeVal(l, tx)(characterId, inventoryType)
			if err != nil {
				return err
			}

			existingItems, err := item.GetByItemId(l, tx)(ii.Id(), itemId)
			if err != nil {
				return err
			}
			// breaks for a rechargeable item.
			if len(existingItems) > 0 {
				index := 0
				for runningQuantity > 0 {
					if index < len(existingItems) {
						i := existingItems[index]
						oldQuantity := uint32(1)
						if val, ok := i.(*item.ItemModel); ok {
							oldQuantity = val.Quantity()
						}
						if oldQuantity > runningQuantity {
							newQuantity := oldQuantity - runningQuantity
							err := item.UpdateQuantity(l, tx)(i.Id(), newQuantity)
							if err != nil {
								l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
								return err
							} else {
								events = append(events, adjustment{mode: 1, itemId: itemId, inventoryType: inventoryType, quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
							}
							runningQuantity = 0
							break
						} else {
							runningQuantity = runningQuantity - oldQuantity
							err := item.RemoveItem(l, tx)(characterId, i.Id())
							if err != nil {
								l.WithError(err).Errorf("Removing quantity %d of item %d.", oldQuantity, i.Id())
								return err
							} else {
								events = append(events, adjustment{mode: 3, itemId: itemId, inventoryType: inventoryType, quantity: oldQuantity, slot: i.Slot(), oldSlot: 0})
							}
						}
						index++
					} else {
						break
					}
				}
			}
			return nil
		})
		if txError != nil {
			return txError
		}

		for _, e := range events {
			emitInventoryModificationEvent(l, span)(characterId, true, e.Mode(), e.ItemId(), e.InventoryType(), e.Quantity(), e.Slot(), e.OldSlot())
		}
		return nil
	}
}
