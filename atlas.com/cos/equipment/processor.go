package equipment

import (
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
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

var characterCreationItems = []uint32{
	1302000, 1312004, 1322005, 1442079, // weapons
	1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167, // bottom
	1060002, 1060006, 1061002, 1061008, 1062115, // top
	1072001, 1072005, 1072037, 1072038, 1072383, // shoes
	30000, 30010, 30020, 30030, 31000, 31040, 31050, // hair
	20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100, //face
}

func CreateForCharacter(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, itemId uint32, equipmentId uint32, characterCreation bool) (*Model, error) {
	return func(characterId uint32, itemId uint32, equipmentId uint32, characterCreation bool) (*Model, error) {
		if characterCreation {
			if invalidCharacterCreationItem(itemId) {
				l.Errorf("Received a request to create an item %d for character %d which is not valid for character creation. This is usually a hack.")
				return nil, errors.New("not valid item for character creation")
			}
		}

		nextOpen, err := getNextFreeEquipmentSlot(db, characterId)
		if err != nil {
			nextOpen = 1
		}

		eq, err := create(db, characterId, equipmentId, nextOpen)
		if err != nil {
			l.Errorf("Persisting equipment %d association for character %d in Slot %d.", equipmentId, characterId, nextOpen)
			return nil, err
		}
		return eq, nil
	}
}

func GetEquipmentForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		return getEquipmentForCharacter(db, characterId)
	}
}

func GetEquippedItemForCharacterBySlot(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, slot int16) (*Model, error) {
	return func(characterId uint32, slot int16) (*Model, error) {
		return getEquipmentForCharacterBySlot(db, characterId, slot)
	}
}

func EquipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32) (*InventoryAdjustment, error) {
	return func(characterId uint32, equipmentId uint32) (*InventoryAdjustment, error) {
		l.Debugf("Received request to equip %d for character %d.", equipmentId, characterId)
		e, err := getByEquipmentId(db, equipmentId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
			return nil, err
		}

		ea, err := statistics.GetEquipmentStatistics(l, span)(e.EquipmentId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
			return nil, err
		}

		l.Debugf("Equipment %d is item %d for character %d.", equipmentId, ea.ItemId(), characterId)

		slots, err := item.GetEquipmentSlotDestination(l, span)(ea.ItemId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve destination slots for item %d.", ea.ItemId())
			return nil, err
		} else if len(slots) <= 0 {
			l.Errorf("Unable to retrieve destination slots for item %d. %s.", ea.ItemId())
			return nil, err
		}
		slot := slots[0]
		l.Debugf("Equipment %d to be equipped in slot %d for character %d.", equipmentId, slot, characterId)

		temporarySlot := int16(math.MinInt16)

		existingSlot := int16(1)
		err = db.Transaction(func(tx *gorm.DB) error {
			if equip, err := getEquipmentForCharacterBySlot(tx, characterId, slot); err == nil && equip.EquipmentId() != 0 {
				l.Debugf("Equipment %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", equip.EquipmentId(), slot, temporarySlot, characterId)
				_ = updateSlot(tx, equip.EquipmentId(), temporarySlot)
			}

			if equip, err := getByEquipmentId(tx, equipmentId); err == nil {
				existingSlot = equip.Slot()
			} else {
				val, err := getNextFreeEquipmentSlot(tx, characterId)
				if err != nil {
				}
				existingSlot = val
			}
			err = updateSlot(tx, equipmentId, slot)
			if err != nil {
				return err
			}
			l.Debugf("Moved item %d from slot %d to %d for character %d.", ea.ItemId(), existingSlot, slot, characterId)

			if equip, err := getEquipmentForCharacterBySlot(tx, characterId, temporarySlot); err == nil && equip.EquipmentId() != 0 {
				err := updateSlot(tx, equip.EquipmentId(), existingSlot)
				if err != nil {
					return err
				}
				l.Debugf("Moved item from temporary location %d to slot %d for character %d.", temporarySlot, existingSlot, characterId)
			}
			return nil
		})
		if err != nil {
			l.WithError(err).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
			return nil, err
		}
		emitItemEquipped(l, span)(characterId)
		return &InventoryAdjustment{
			mode:          2,
			itemId:        ea.ItemId(),
			inventoryType: 1,
			quantity:      1,
			slot:          slot,
			oldSlot:       existingSlot,
		}, nil
	}
}

func UnequipItemForCharacter(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, equipmentId uint32, oldSlot int16) (*InventoryAdjustment, error) {
	return func(characterId uint32, equipmentId uint32, oldSlot int16) (*InventoryAdjustment, error) {
		l.Debugf("Received request to unequip %d for character %d.", equipmentId, characterId)

		itemId := uint32(0)
		newSlot := int16(0)
		txErr := db.Transaction(func(tx *gorm.DB) error {
			ea, err := statistics.GetEquipmentStatistics(l, span)(equipmentId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
				return err
			}
			itemId = ea.ItemId()

			newSlot, err = getNextFreeEquipmentSlot(tx, characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to get next free equipment slot")
				return err
			}

			err = updateSlot(tx, equipmentId, newSlot)
			if err != nil {
				return err
			}

			l.Debugf("Unequipped %d for character %d and place it in slot %d, from %d.", equipmentId, characterId, newSlot, oldSlot)
			return nil
		})
		if txErr != nil {
			l.WithError(txErr).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
			return nil, txErr
		}
		emitItemUnequipped(l, span)(characterId)
		return &InventoryAdjustment{
			mode:          2,
			itemId:        itemId,
			inventoryType: 1,
			quantity:      1,
			slot:          newSlot,
			oldSlot:       oldSlot,
		}, nil
	}
}

func GetEquipmentById(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		return getById(db, id)
	}
}

func invalidCharacterCreationItem(itemId uint32) bool {
	for _, v := range characterCreationItems {
		if itemId == v {
			return false
		}
	}
	return true
}

func GainItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(characterId uint32, itemId uint32, equipmentId uint32) (*InventoryAdjustment, error) {
	return func(characterId uint32, itemId uint32, equipmentId uint32) (*InventoryAdjustment, error) {
		//TODO verify inventory space
		e, err := CreateForCharacter(l, db)(characterId, itemId, equipmentId, false)
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
			return nil, err
		}
		return &InventoryAdjustment{
			mode:          0,
			itemId:        itemId,
			inventoryType: 1,
			quantity:      1,
			slot:          e.Slot(),
			oldSlot:       0,
		}, nil
	}
}

func DropEquippedItem(l logrus.FieldLogger, db *gorm.DB, _ opentracing.Span) func(worldId byte, channelId byte, characterId uint32, slot int16) (uint32, error) {
	return func(worldId byte, channelId byte, characterId uint32, slot int16) (uint32, error) {
		l.Debugf("Character %d dropping equipment in slot %d.", characterId, slot)
		e, err := GetEquippedItemForCharacterBySlot(l, db)(characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return 0, err
		}

		err = RemoveItem(l, db)(characterId, e.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to remove item %d.", e.Id())
			return 0, err
		}
		return e.EquipmentId(), nil
	}
}

func DropEquipment(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, slot int16) (uint32, error) {
	return func(worldId byte, channelId byte, characterId uint32, slot int16) (uint32, error) {
		l.Debugf("Character %d dropping equipment in slot %d.", characterId, slot)
		e, err := GetEquippedItemForCharacterBySlot(l, db)(characterId, slot)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve item in slot %d being dropped.", slot)
			return 0, err
		}

		err = RemoveItem(l, db)(characterId, e.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to remove item %d.", e.Id())
			return 0, err
		}

		emitItemUnequipped(l, span)(characterId)
		return e.EquipmentId(), nil
	}
}

func RemoveItem(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, id uint32) error {
	return func(characterId uint32, id uint32) error {
		return remove(db, characterId, id)
	}
}

func MoveItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, source int16, destination int16) (*InventoryAdjustment, error) {
	return func(characterId uint32, source int16, destination int16) (*InventoryAdjustment, error) {
		itemId := uint32(0)
		txError := db.Transaction(func(tx *gorm.DB) error {
			equip, err := getEquipmentForCharacterBySlot(tx, characterId, source)
			if err != nil || equip.Id() == 0 {
				l.Warnf("Item movement requested, but no equipment for character %d in slot %d.", characterId, source)
				return err
			}

			ea, err := statistics.GetEquipmentStatistics(l, span)(equip.EquipmentId())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve equipment %d.", equip.EquipmentId())
				return err
			}
			itemId = ea.ItemId()

			temporarySlot := int16(math.MinInt16)
			otherEquip, err := getEquipmentForCharacterBySlot(tx, characterId, destination)
			if err == nil && otherEquip.Id() != 0 {
				l.Debugf("Equipment %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", otherEquip.Id(), destination, temporarySlot, characterId)
				_ = updateSlot(tx, otherEquip.EquipmentId(), temporarySlot)
			}

			err = updateSlot(tx, equip.EquipmentId(), destination)
			if err != nil {
				return err
			}
			l.Debugf("Moved item %d from slot %d to %d for character %d.", equip.Id(), source, destination, characterId)

			if otherEquip != nil {
				err = updateSlot(tx, otherEquip.EquipmentId(), source)
				if err != nil {
					return err
				}
				l.Debugf("Moved item %d from slot %d to %d for character %d.", otherEquip.Id(), temporarySlot, source, characterId)
			}

			return nil
		})
		return &InventoryAdjustment{
			mode:          2,
			itemId:        itemId,
			inventoryType: 1,
			quantity:      1,
			slot:          destination,
			oldSlot:       source,
		}, txError
	}
}

func GetItemIdForEquipment(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) (uint32, error) {
	return func(equipmentId uint32) (uint32, error) {
		ea, err := statistics.GetEquipmentStatistics(l, span)(equipmentId)
		if err != nil {
			return 0, err
		}
		return ea.ItemId(), nil
	}
}
