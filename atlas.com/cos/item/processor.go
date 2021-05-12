package item

import (
	"atlas-cos/kafka/producers"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

func GetItemsForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, itemId uint32) []*Model {
	return func(characterId uint32, inventoryType int8, itemId uint32) []*Model {

		items, err := getItemsForCharacter(db, characterId, inventoryType, itemId)
		if err != nil {
			return make([]*Model, 0)
		}
		return items
	}
}

func GetItemForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, slot int16) (*Model, error) {
	return func(characterId uint32, inventoryType int8, slot int16) (*Model, error) {
		return getItemForCharacter(db, characterId, inventoryType, slot)
	}
}

func GetForCharacterByInventory(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) ([]*Model, error) {
	return func(characterId uint32, inventoryType int8) ([]*Model, error) {
		return getForCharacterByInventory(db, characterId, inventoryType)
	}
}

func UpdateItemQuantity(_ logrus.FieldLogger, db *gorm.DB) func(id uint32, quantity uint32) error {
	return func(id uint32, quantity uint32) error {
		return update(db, id, SetQuantity(quantity))
	}
}

func CreateItemForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (*Model, error) {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (*Model, error) {
		slot, err := getNextFreeEquipmentSlot(db, characterId, inventoryType)
		if err != nil {
			return nil, err
		}
		return createItemForCharacter(db, characterId, inventoryType, itemId, quantity, slot)
	}
}

func GetItemById(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		return getById(db, id)
	}
}

func RemoveItem(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) error {
	return func(id uint32) error {
		return remove(db, id)
	}
}

func maxInSlot() uint32 {
	return 200
}

func GainItem(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, it int8, itemId uint32, quantity uint32) error {
	return func(characterId uint32, it int8, itemId uint32, quantity uint32) error {
		//TODO verify inventory space

		slotMax := maxInSlot()
		runningQuantity := quantity

		existingItems := GetItemsForCharacter(l, db)(characterId, it, itemId)
		// breaks for a rechargeable item.
		if len(existingItems) > 0 {
			index := 0
			for runningQuantity > 0 {
				if index < len(existingItems) {
					i := existingItems[index]
					oldQuantity := i.Quantity()
					if oldQuantity < slotMax {
						newQuantity := uint32(math.Min(float64(oldQuantity+runningQuantity), float64(slotMax)))
						runningQuantity = runningQuantity - (newQuantity - oldQuantity)
						err := UpdateItemQuantity(l, db)(i.Id(), newQuantity)
						if err != nil {
							l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
						} else {
							producers.InventoryModificationReservation(l, context.Background()).
								Emit(characterId, true, 1, itemId, i.InventoryType(), newQuantity, i.Slot(), 0)
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
			i, err := CreateItemForCharacter(l, db)(characterId, it, itemId, newQuantity)
			if err != nil {
				l.WithError(err).Errorf("Unable to create item %d that character %d picked up.", itemId, characterId)
				return err
			}
			producers.InventoryModificationReservation(l, context.Background()).
				Emit(characterId, true, 0, itemId, it, newQuantity, i.Slot(), 0)
		}

		return nil
	}
}

func LoseItem(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, it int8, itemId uint32, quantity int32) error {
	return func(characterId uint32, it int8, itemId uint32, quantity int32) error {
		//TODO verify inventory space

		runningQuantity := uint32(quantity * -1)

		existingItems := GetItemsForCharacter(l, db)(characterId, it, itemId)
		// breaks for a rechargeable item.
		if len(existingItems) > 0 {
			index := 0
			for runningQuantity > 0 {
				if index < len(existingItems) {
					i := existingItems[index]
					oldQuantity := i.Quantity()
					if oldQuantity > runningQuantity {
						newQuantity := oldQuantity - runningQuantity
						err := UpdateItemQuantity(l, db)(i.Id(), newQuantity)
						if err != nil {
							l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
						} else {
							producers.InventoryModificationReservation(l, context.Background()).
								Emit(characterId, true, 1, itemId, i.InventoryType(), newQuantity, i.Slot(), 0)
						}
						runningQuantity = 0
						break
					} else {
						runningQuantity = runningQuantity - oldQuantity
						err := RemoveItem(l, db)(i.Id())
						if err != nil {
							l.WithError(err).Errorf("Removing quantity %d of item %d.", oldQuantity, i.Id())
						} else {
							producers.InventoryModificationReservation(l, context.Background()).
								Emit(characterId, true, 3, itemId, i.InventoryType(), oldQuantity, i.Slot(), 0)
						}
					}
					index++
				} else {
					break
				}
			}
		}
		return nil
	}
}

func DropItem(l logrus.FieldLogger, db *gorm.DB) func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) (uint32, error) {
	return func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) (uint32, error) {
		l.Debugf("Character %d dropping %d item in inventory %d slot %d.", characterId, quantity, inventoryType, slot)
		i, err := GetItemForCharacter(l, db)(characterId, inventoryType, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return 0, err
		}

		if i.Quantity() <= uint32(quantity) {
			err := RemoveItem(l, db)(i.Id())
			if err != nil {
				l.WithError(err).Errorf("Could not remove item %d from character %d inventory.", i.Id(), characterId)
				return 0, err
			}
			producers.InventoryModificationReservation(l, context.Background()).
				Emit(characterId, true, 3, i.ItemId(), i.InventoryType(), uint32(quantity), i.Slot(), 0)
		} else {
			newQuantity := i.Quantity() - uint32(quantity)
			err := UpdateItemQuantity(l, db)(i.Id(), newQuantity)
			if err != nil {
				l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
				return 0, err
			} else {
				producers.InventoryModificationReservation(l, context.Background()).
					Emit(characterId, true, 1, i.ItemId(), i.InventoryType(), newQuantity, i.Slot(), 0)
			}
		}

		return i.ItemId(), nil
	}
}
