package drop

import (
	"atlas-cos/character"
	"atlas-cos/inventory"
	"atlas-cos/inventory/item"
	"atlas-cos/model"
	"atlas-cos/party"
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

func AttemptPickup(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, dropId uint32) {
	return func(characterId uint32, dropId uint32) {
		c, err := character.GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Attempting to pick up %d for character %d.", dropId, characterId)
			return
		}
		l.Debugf("Character %d attempting to pickup drop %d.", characterId, dropId)
		d, err := GetById(l, span)(dropId)
		if err != nil {
			l.WithError(err).Errorf("Attempting to pick up %d for character %d.", dropId, characterId)
			return
		}
		attemptPickup(l, db, span)(c, d)
	}
}

func byIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32) model.Provider[Model] {
	return func(dropId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(dropId), makeDrop)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32) (Model, error) {
	return func(dropId uint32) (Model, error) {
		return byIdModelProvider(l, span)(dropId)()
	}
}

func makeDrop(dc requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.Atoi(dc.Id)
	if err != nil {
		return Model{}, nil
	}
	attr := dc.Attributes
	return Model{
		id:            uint32(id),
		itemId:        attr.ItemId,
		equipmentId:   attr.EquipmentId,
		quantity:      attr.Quantity,
		meso:          attr.Meso,
		dropTime:      attr.DropTime,
		dropType:      attr.DropType,
		ownerId:       attr.OwnerId,
		characterDrop: attr.CharacterDrop,
	}, nil
}

func attemptPickup(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c character.Model, d Model) {
	return func(c character.Model, d Model) {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		elapsed := now - int64(d.DropTime())
		if elapsed < 400 {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop has not yet met minimum time. Time elapsed %d", c.Id(), d.Id(), elapsed)
			l.Debugf("Now %d, DropTime %d, Elapsed %d.", now, d.DropTime(), elapsed)
			emitCancelReservation(l, span)(d.Id(), c.Id())
			return
		}

		if !canBePickedBy(l, span)(c, d) {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop cannot be picked up by character.", c.Id(), d.Id())
			emitCancelReservation(l, span)(d.Id(), c.Id())
			return
		}

		if isOwnerLockedMap(c.MapId()) && d.CharacterDrop() && d.OwnerId() != c.Id() {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop is not owned by this character, in a owner locked map.", c.Id(), d.Id())
			emitCancelReservation(l, span)(d.Id(), c.Id())
			// emit item unavailable.
			return
		}

		if d.ItemId() == 4031865 || d.ItemId() == 4031866 {
			l.Debugf("Picking up NX item %d for character %d.", d.ItemId(), c.Id())
			pickupNX(l, span)(c, d)
		} else if d.Meso() > 0 {
			l.Debugf("Picking up meso drop for character %d.", c.Id())
			pickupMeso(l, span)(c, d)
		} else if consumeOnPickup() {
		} else {

			if !needsQuestItem() {
				l.Debugf("Cancelling drop for character %d, drop %d, the character does not need this quest item.", c.Id(), d.Id())
				emitCancelReservation(l, span)(d.Id(), c.Id())
				// emit item unavailable.
				return
			}

			if !hasInventorySpace(l, db)(c.Id(), d) {
				l.Debugf("Cancelling drop for character %d, drop %d, the character does not have inventory space.", c.Id(), d.Id())
				emitCancelReservation(l, span)(d.Id(), c.Id())
				emitInventoryFullEvent(l, span)(c.Id())
				// emit inventory full.
				return
			}

			if scriptedItem(d.ItemId()) {
				// TODO handle scripted item
			}

			if val, ok := inventory.GetInventoryType(d.ItemId()); ok {
				if val == inventory.TypeValueEquip {
					l.Debugf("Picking up equipment %d for character %d.", d.EquipmentId(), c.Id())
					pickupEquip(l, db, span)(c, d)
				} else {
					l.Debugf("Picking up item %d for character %d.", d.ItemId(), c.Id())
					pickupItem(l, db, span)(c, val, d)
				}
				// TODO update ariant score if 4031868
				emitPickedUpItemEvent(l, span)(c.Id(), d.ItemId(), d.Quantity())
			}
		}
		emitPickedUpDropCommand(l, span)(c.Id(), d.Id())
	}
}

func canBePickedBy(l logrus.FieldLogger, span opentracing.Span) func(c character.Model, d Model) bool {
	return func(c character.Model, d Model) bool {
		if d.OwnerId() <= 0 || d.FFADrop() {
			return true
		}

		ownerParty, err := party.ForCharacter(l, span)(d.OwnerId())
		if err != nil {
			if c.Id() == d.OwnerId() {
				return true
			}
		} else {
			characterParty, err := party.ForCharacter(l, span)(c.Id())
			if err == nil && ownerParty.Id() == characterParty.Id() {
				return true
			} else if c.Id() == d.OwnerId() {
				return true
			}
		}
		return d.ExpiredOwnershipTime()
	}
}

func isOwnerLockedMap(mapId uint32) bool {
	return mapId > 209000000 && mapId < 209000016 || mapId >= 990000500 && mapId <= 990000502
}

func pickupNX(l logrus.FieldLogger, span opentracing.Span) func(c character.Model, d Model) {
	return func(c character.Model, d Model) {
		if d.ItemId() == 4031865 {
			emitPickedUpNxEvent(l, span)(c.Id(), 100)
		} else {
			emitPickedUpNxEvent(l, span)(c.Id(), 250)
		}
	}
}

func pickupMeso(l logrus.FieldLogger, span opentracing.Span) func(c character.Model, d Model) {
	return func(c character.Model, d Model) {
		emitMesoAdjustment(l, span)(c.Id(), int32(d.Meso()))
	}
}

func consumeOnPickup() bool {
	return false
}

func needsQuestItem() bool {
	return true
}

func hasInventorySpace(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, d Model) bool {
	return func(characterId uint32, d Model) bool {
		//TODO checking inventory space, and adding items should become an atomic action
		if val, ok := inventory.GetInventoryType(d.ItemId()); ok {
			i, err := inventory.GetInventoryByTypeVal(l, db)(characterId, val)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", characterId)
				return false
			}

			if val == inventory.TypeValueEquip {
				return hasEquipInventorySpace(l)(characterId, i, d.ItemId())
			} else {
				return hasItemInventorySpace(l)(characterId, i, d.ItemId(), d.Quantity())
			}
		}
		return false
	}
}

func hasItemInventorySpace(l logrus.FieldLogger) func(characterId uint32, inventory inventory.Model, itemId uint32, quantity uint32) bool {
	return func(characterId uint32, inventory inventory.Model, itemId uint32, quantity uint32) bool {
		l.Debugf("Checking inventory capacity for item %d, quantity %d for character %d.", itemId, quantity, characterId)
		slotMax := maxInSlot()
		runningQuantity := quantity
		existingItems := inventory.Items()

		// breaks for a rechargeable item.
		usedSlots := uint32(len(existingItems))

		l.Debugf("Character %d has %d slots already occupied.", characterId, usedSlots)

		if len(existingItems) > 0 {
			index := 0
			for runningQuantity > 0 {
				if index >= len(existingItems) {
					break
				}

				i := existingItems[index]
				if i.ItemId() == itemId {
					oldQuantity := uint32(0)
					if val, ok := i.(item.ItemModel); ok {
						oldQuantity = val.Quantity()
					}
					if oldQuantity < slotMax {
						newQuantity := uint32(math.Min(float64(oldQuantity+runningQuantity), float64(slotMax)))
						runningQuantity = runningQuantity - (newQuantity - oldQuantity)
					}
				}
				index++
			}
		}

		newSlots := uint32(0)
		for runningQuantity > 0 {
			newQuantity := uint32(math.Min(float64(runningQuantity), float64(slotMax)))
			runningQuantity = runningQuantity - newQuantity
			newSlots += 1
		}
		l.Debugf("Character %d would need to consume %d additional slot(s) to pick up %d %d.", characterId, newSlots, quantity, itemId)

		if usedSlots+newSlots > inventory.Capacity() {
			l.Debugf("Unable to pick up item %d because character %d inventory full.", itemId, characterId)
			return false
		}
		return true
	}
}

func hasEquipInventorySpace(l logrus.FieldLogger) func(characterId uint32, inventory inventory.Model, itemId uint32) bool {
	return func(characterId uint32, inventory inventory.Model, itemId uint32) bool {
		l.Debugf("Checking inventory capacity for equip %d for character %d.", itemId, characterId)
		count := uint32(0)
		for _, equip := range inventory.Items() {
			if equip.Slot() >= 0 {
				count += 1
			}
		}
		if count+1 > inventory.Capacity() {
			l.Debugf("Unable to pick up equip %d because character %d inventory full.", itemId, characterId)
			return false
		}
		return true
	}
}

func scriptedItem(itemId uint32) bool {
	return itemId/10000 == 243
}

func pickupEquip(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c character.Model, d Model) {
	return func(c character.Model, d Model) {
		err := inventory.GainEquipment(l, db, span)(c.Id(), d.ItemId(), d.EquipmentId())
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d that character %d picked up.", d.ItemId(), c.Id())
			return
		}
	}
}

func pickupItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c character.Model, it int8, d Model) {
	return func(c character.Model, it int8, d Model) {
		err := inventory.GainItem(l, db, span)(c.Id(), it, d.ItemId(), d.Quantity())
		if err != nil {
			l.WithError(err).Errorf("Unable to gain item %d that character %d picked up.", d.ItemId(), c.Id())
			return
		}
	}
}

func maxInSlot() uint32 {
	return 200
}
