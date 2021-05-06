package equipment

import (
	"atlas-cos/kafka/producers"
	"atlas-cos/rest/requests"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
)

var characterCreationItems = []uint32{
	1302000, 1312004, 1322005, 1442079, // weapons
	1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167, // bottom
	1060002, 1060006, 1061002, 1061008, 1062115, // top
	1072001, 1072005, 1072037, 1072038, 1072383, // shoes
	30000, 30010, 30020, 30030, 31000, 31040, 31050, // hair
	20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100, //face
}

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) CreateForCharacter(characterId uint32, itemId uint32, characterCreation bool) (*Model, error) {
	if characterCreation {
		if invalidCharacterCreationItem(itemId) {
			p.l.Errorf("Received a request to create an item %d for character %d which is not valid for character creation. This is usually a hack.")
			return nil, errors.New("not valid item for character creation")
		}
	}

	nextOpen, err := GetNextFreeEquipmentSlot(p.db, characterId)
	if err != nil {
		nextOpen = 1
	}

	ro, err := requests.EquipmentRegistry().Create(itemId)
	if err != nil {
		p.l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
		return nil, err
	}
	eid, err := strconv.Atoi(ro.Data.Id)
	if err != nil {
		p.l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.")
		return nil, err
	}
	eq, err := Create(p.db, characterId, uint32(eid), nextOpen)
	if err != nil {
		p.l.Errorf("Persisting equipment %d association for character %d in Slot %d.", eid, characterId, nextOpen)
		return nil, err
	}
	return eq, nil
}

func (p processor) GetEquipmentForCharacter(characterId uint32) ([]*Model, error) {
	return GetEquipmentForCharacter(p.db, characterId)
}

func (p processor) GetEquippedItemForCharacterBySlot(characterId uint32, slot int16) (*Model, error) {
	return GetEquipmentForCharacterBySlot(p.db, characterId, slot)
}

func (p processor) CreateAndEquip(characterId uint32, items ...uint32) {
	for _, item := range items {
		p.createAndEquip(characterId, item)
	}
}

func (p processor) createAndEquip(characterId uint32, itemId uint32) {
	if equipment, err := p.CreateForCharacter(characterId, itemId, true); err == nil {
		p.EquipItemForCharacter(characterId, equipment.EquipmentId())
	} else {
		p.l.Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
	}
}

func (p processor) EquipItemForCharacter(characterId uint32, equipmentId uint32) {
	p.l.Debugf("Received request to equip %d for character %d.", equipmentId, characterId)
	e, err := GetByEquipmentId(p.db, equipmentId)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
		return
	}

	ea, err := requests.EquipmentRegistry().GetById(e.EquipmentId())
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
		return
	}

	itemId := ea.Data.Attributes.ItemId
	p.l.Debugf("Equipment %d is item %d for character %d.", equipmentId, itemId, characterId)

	slots, err := p.getEquipmentSlotDestination(itemId)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve destination slots for item %d.", itemId)
		return
	} else if len(slots) <= 0 {
		p.l.Errorf("Unable to retrieve destination slots for item %d. %s.", itemId)
		return
	}
	slot := slots[0]
	p.l.Debugf("Equipment %d to be equipped in slot %d for character %d.", equipmentId, slot, characterId)

	temporarySlot := int16(math.MinInt16)

	events := make([]func(), 0)

	err = p.db.Transaction(func(tx *gorm.DB) error {
		if equip, err := GetEquipmentForCharacterBySlot(tx, characterId, slot); err == nil && equip.EquipmentId() != 0 {
			p.l.Debugf("Equipment %d already exists in slot %d, that item will be moved temporarily to %d for character %d.", equip.EquipmentId(), slot, temporarySlot, characterId)
			_ = UpdateSlot(tx, equip.EquipmentId(), temporarySlot)
		}

		currentSlot := int16(1)
		if equip, err := GetByEquipmentId(tx, equipmentId); err == nil {
			currentSlot = equip.Slot()
		} else {
			val, err := GetNextFreeEquipmentSlot(tx, characterId)
			if err != nil {
			}
			currentSlot = val
		}
		err = UpdateSlot(tx, equipmentId, slot)
		if err != nil {
			return err
		}
		p.l.Debugf("Moved item %d from slot %d to %d for character %d.", itemId, currentSlot, slot, characterId)

		events = append(events, func() {
			producers.InventoryModificationReservation(p.l, context.Background()).Emit(characterId, true, 2, ea.Data.Attributes.ItemId, 1, 1, slot, currentSlot)
		})

		if equip, err := GetEquipmentForCharacterBySlot(tx, characterId, temporarySlot); err == nil && equip.EquipmentId() != 0 {
			err := UpdateSlot(tx, equip.EquipmentId(), currentSlot)
			if err != nil {
				return err
			}
			p.l.Debugf("Moved item from temporary location %d to slot %d for character %d.", temporarySlot, currentSlot, characterId)
		}
		events = append(events, func() {
			producers.CharacterEquippedItem(p.l)(characterId)
		})
		return nil
	})
	if err != nil {
		p.l.WithError(err).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
		return
	}

	for _, event := range events {
		event()
	}
}

func (p processor) UnequipItemForCharacter(characterId uint32, equipmentId uint32, oldSlot int16) {
	p.l.Debugf("Received request to unequip %d for character %d.", equipmentId, characterId)

	events := make([]func(), 0)

	err := p.db.Transaction(func(tx *gorm.DB) error {
		ea, err := requests.EquipmentRegistry().GetById(equipmentId)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to retrieve equipment %d.", equipmentId)
			return err
		}

		val, err := GetNextFreeEquipmentSlot(tx, characterId)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to get next free equipment slot")
			return err
		}

		err = UpdateSlot(tx, equipmentId, val)
		if err != nil {
			return err
		}

		p.l.Debugf("Unequipped %d for character %d and place it in slot %d, from %d.", equipmentId, characterId, val, oldSlot)
		events = append(events, func() {
			producers.InventoryModificationReservation(p.l, context.Background()).Emit(characterId, true, 2, ea.Data.Attributes.ItemId, 1, 1, val, oldSlot)
		})
		events = append(events, func() {
			producers.CharacterUnEquippedItem(p.l)(characterId)
		})
		return nil
	})
	if err != nil {
		p.l.WithError(err).Errorf("Unable to complete the equipment of item %d for character %d.", equipmentId, characterId)
		return
	}

	for _, event := range events {
		event()
	}
}

func (p processor) getEquipmentSlotDestination(itemId uint32) ([]int16, error) {
	r, err := requests.ItemInformationRegistry().GetEquipmentSlotDestination(itemId)
	if err != nil {
		return nil, err
	}

	var slots = make([]int16, 0)
	for _, data := range r.Data {
		slots = append(slots, data.Attributes.Slot)
	}
	return slots, nil
}

func (p processor) GetEquipmentById(id uint32) (*Model, error) {
	return GetById(p.db, id)
}

func invalidCharacterCreationItem(itemId uint32) bool {
	for _, v := range characterCreationItems {
		if itemId == v {
			return false
		}
	}
	return true
}

func GainItem(l log.FieldLogger, db *gorm.DB) func(characterId uint32, itemId uint32) error {
	return func(characterId uint32, itemId uint32) error {
		//TODO verify inventory space
		e, err := Processor(l, db).CreateForCharacter(characterId, itemId, false)
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d for character %d.", itemId, characterId)
			return err
		}
		producers.InventoryModificationReservation(l, context.Background()).
			Emit(characterId, true, 0, itemId, 1, 1, e.Slot(), 0)
		return nil
	}
}
