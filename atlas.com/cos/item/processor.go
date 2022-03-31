package item

import (
	"atlas-cos/database"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"sort"
)

type InventoryAdjustment struct {
	mode          byte
	itemId        uint32
	inventoryType int8
	quantity      uint32
	slot          int16
	oldSlot       int16
}

func (i InventoryAdjustment) Mode() byte {
	return i.mode
}

func (i InventoryAdjustment) ItemId() uint32 {
	return i.itemId
}

func (i InventoryAdjustment) InventoryType() int8 {
	return i.inventoryType
}

func (i InventoryAdjustment) Quantity() uint32 {
	return i.quantity
}

func (i InventoryAdjustment) Slot() int16 {
	return i.slot
}

func (i InventoryAdjustment) OldSlot() int16 {
	return i.oldSlot
}

func GetItemsForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, itemId uint32) ([]Model, error) {
	return func(characterId uint32, inventoryType int8, itemId uint32) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(getItemsForCharacter(characterId, inventoryType, itemId), makeItem)()
	}
}

func GetItemForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, slot int16) (Model, error) {
	return func(characterId uint32, inventoryType int8, slot int16) (Model, error) {
		return database.ModelProvider[Model, entity](db)(getItemForCharacter(characterId, inventoryType, slot), makeItem)()
	}
}

func UpdateItemQuantity(_ logrus.FieldLogger, db *gorm.DB) func(id uint32, quantity uint32) error {
	return func(id uint32, quantity uint32) error {
		return update(db, id, SetQuantity(quantity))
	}
}

func CreateItemForCharacter(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (Model, error) {
	return func(characterId uint32, inventoryType int8, itemId uint32, quantity uint32) (Model, error) {
		slot, err := getNextFreeEquipmentSlot(l, db)(characterId, inventoryType)
		if err != nil {
			return Model{}, err
		}
		return createItemForCharacter(db, characterId, inventoryType, itemId, quantity, slot)
	}
}

func GetItemsForCharacterByInventory(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) ([]Model, error) {
	return func(characterId uint32, inventoryType int8) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(getItemsForCharacterByInventory(characterId, inventoryType), makeItem)()
	}
}

func getNextFreeEquipmentSlot(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryType int8) (int16, error) {
	return func(characterId uint32, inventoryType int8) (int16, error) {
		items, err := GetItemsForCharacterByInventory(l, db)(characterId, inventoryType)
		if err != nil {
			return 1, err
		}

		l.Debugf("Character %d has %d items in %d inventory.", characterId, len(items), inventoryType)

		if len(items) == 0 {
			l.Debugf("Defaulting item slot choice to 1 for character %d.", characterId)
			return 1, nil
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Slot() < items[j].Slot()
		})
		slot := minFreeSlot(items)

		l.Debugf("Chose slot %d for new item in inventory %d for character %d.", slot, inventoryType, characterId)

		return slot, nil
	}
}

func minFreeSlot(items []Model) int16 {
	slot := int16(1)
	i := 0

	for {
		if i >= len(items) {
			return slot
		} else if slot < items[i].Slot() {
			return slot
		} else if slot == items[i].Slot() {
			slot += 1
			i += 1
		} else if items[i].Slot() <= 0 {
			i += 1
		}
	}
}

func GetItemById(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) (Model, error) {
	return func(id uint32) (Model, error) {
		return database.ModelProvider[Model, entity](db)(getById(id), makeItem)()
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

func GainItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, it int8, itemId uint32, quantity uint32) ([]InventoryAdjustment, error) {
	return func(characterId uint32, it int8, itemId uint32, quantity uint32) ([]InventoryAdjustment, error) {
		//TODO verify inventory space

		slotMax := maxInSlot()
		runningQuantity := quantity

		var events = make([]InventoryAdjustment, 0)

		txError := db.Transaction(func(tx *gorm.DB) error {
			existingItems, err := GetItemsForCharacter(l, tx)(characterId, it, itemId)
			if err != nil {
				return err
			}
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
							err := UpdateItemQuantity(l, tx)(i.Id(), newQuantity)
							if err != nil {
								l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
							} else {
								events = append(events, InventoryAdjustment{mode: 1, itemId: itemId, inventoryType: i.InventoryType(), quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
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
				i, err := CreateItemForCharacter(l, tx)(characterId, it, itemId, newQuantity)
				if err != nil {
					l.WithError(err).Errorf("Unable to create item %d that character %d picked up.", itemId, characterId)
					return err
				}
				events = append(events, InventoryAdjustment{mode: 0, itemId: itemId, inventoryType: it, quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
			}
			return nil
		})
		return events, txError
	}
}

func LoseItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, it int8, itemId uint32, quantity int32) ([]InventoryAdjustment, error) {
	return func(characterId uint32, it int8, itemId uint32, quantity int32) ([]InventoryAdjustment, error) {
		//TODO verify inventory space

		runningQuantity := uint32(quantity * -1)

		var events = make([]InventoryAdjustment, 0)

		txError := db.Transaction(func(tx *gorm.DB) error {
			existingItems, err := GetItemsForCharacter(l, tx)(characterId, it, itemId)
			if err != nil {
				return err
			}
			// breaks for a rechargeable item.
			if len(existingItems) > 0 {
				index := 0
				for runningQuantity > 0 {
					if index < len(existingItems) {
						i := existingItems[index]
						oldQuantity := i.Quantity()
						if oldQuantity > runningQuantity {
							newQuantity := oldQuantity - runningQuantity
							err := UpdateItemQuantity(l, tx)(i.Id(), newQuantity)
							if err != nil {
								l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
								return err
							} else {
								events = append(events, InventoryAdjustment{mode: 1, itemId: itemId, inventoryType: i.InventoryType(), quantity: newQuantity, slot: i.Slot(), oldSlot: 0})
							}
							runningQuantity = 0
							break
						} else {
							runningQuantity = runningQuantity - oldQuantity
							err := RemoveItem(l, tx)(i.Id())
							if err != nil {
								l.WithError(err).Errorf("Removing quantity %d of item %d.", oldQuantity, i.Id())
								return err
							} else {
								events = append(events, InventoryAdjustment{mode: 3, itemId: itemId, inventoryType: i.InventoryType(), quantity: oldQuantity, slot: i.Slot(), oldSlot: 0})
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
		return events, txError
	}
}

func DropItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) (*InventoryAdjustment, error) {
	return func(worldId byte, channelId byte, characterId uint32, inventoryType int8, slot int16, quantity int16) (*InventoryAdjustment, error) {
		l.Debugf("Character %d dropping %d item in inventory %d slot %d.", characterId, quantity, inventoryType, slot)
		i, err := GetItemForCharacter(l, db)(characterId, inventoryType, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return nil, err
		}

		if i.Quantity() <= uint32(quantity) {
			err = RemoveItem(l, db)(i.Id())
			if err != nil {
				l.WithError(err).Errorf("Could not remove item %d from character %d inventory.", i.Id(), characterId)
				return nil, err
			}
			return &InventoryAdjustment{
				mode:          3,
				itemId:        i.ItemId(),
				inventoryType: inventoryType,
				quantity:      uint32(quantity),
				slot:          slot,
				oldSlot:       0,
			}, nil
		} else {
			newQuantity := i.Quantity() - uint32(quantity)
			err = UpdateItemQuantity(l, db)(i.Id(), newQuantity)
			if err != nil {
				l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
				return nil, err
			}
			return &InventoryAdjustment{
				mode:          1,
				itemId:        i.ItemId(),
				inventoryType: inventoryType,
				quantity:      newQuantity,
				slot:          slot,
				oldSlot:       0,
			}, nil
		}
	}
}

func MoveItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, inventoryType int8, source int16, destination int16) (*InventoryAdjustment, error) {
	return func(characterId uint32, inventoryType int8, source int16, destination int16) (*InventoryAdjustment, error) {
		itemId := uint32(0)
		quantity := uint32(0)
		txError := db.Transaction(func(tx *gorm.DB) error {
			item, err := GetItemForCharacter(l, tx)(characterId, inventoryType, source)
			if err != nil || item.Id() == 0 {
				l.Warnf("Item movement requested, but no equipment for character %d in slot %d.", characterId, source)
				return err
			}
			itemId = item.ItemId()
			quantity = item.Quantity()

			temporarySlot := int16(math.MinInt16)
			otherItem, err := GetItemForCharacter(l, tx)(characterId, inventoryType, destination)
			if err == nil && otherItem.Id() != 0 {
				l.Debugf("Item %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", otherItem.Id(), destination, temporarySlot, characterId)
				err = update(tx, otherItem.Id(), SetSlot(temporarySlot))
				if err != nil {
					l.WithError(err).Errorf("Unable to move other item %d to temporary slot %d.", otherItem.Id(), temporarySlot)
					return err
				}
			}

			err = update(tx, item.Id(), SetSlot(destination))
			if err != nil {
				return err
			}
			l.Debugf("Moved item %d from slot %d to %d for character %d.", item.Id(), source, destination, characterId)

			if otherItem.Id() != 0 {
				err = update(tx, otherItem.Id(), SetSlot(source))
				if err != nil {
					l.WithError(err).Errorf("Unable to move other item %d to resulting slot %d.", otherItem.Id(), source)
					return err
				}
				l.Debugf("Moved item %d from slot %d to %d for character %d.", otherItem.Id(), temporarySlot, source, characterId)
			}
			return nil
		})
		return &InventoryAdjustment{
			mode:          2,
			itemId:        itemId,
			inventoryType: inventoryType,
			quantity:      quantity,
			slot:          destination,
			oldSlot:       source,
		}, txError
	}
}

func GetEquipmentSlotDestination(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) ([]int16, error) {
	return func(itemId uint32) ([]int16, error) {
		r, err := requestEquipmentSlotDestination(itemId)(l, span)
		if err != nil {
			return nil, err
		}

		var slots = make([]int16, 0)
		for _, data := range r.DataList() {
			attr := data.Attributes
			slots = append(slots, attr.Slot)
		}
		return slots, nil
	}
}
