package drop

import (
	"atlas-cos/character"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/kafka/producers"
	"atlas-cos/party"
	"atlas-cos/rest/attributes"
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

func GetById(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32) (*Model, error) {
	return func(dropId uint32) (*Model, error) {
		dc, err := requestById(l, span)(dropId)
		if err != nil {
			return nil, err
		}
		d := makeDrop(dc.Data)
		return d, nil
	}
}

func makeDrop(dc attributes.DropData) *Model {
	id, err := strconv.Atoi(dc.Id)
	if err != nil {
		return nil
	}
	return &Model{
		id:            uint32(id),
		itemId:        dc.Attributes.ItemId,
		equipmentId:   dc.Attributes.EquipmentId,
		quantity:      dc.Attributes.Quantity,
		meso:          dc.Attributes.Meso,
		dropTime:      dc.Attributes.DropTime,
		dropType:      dc.Attributes.DropType,
		ownerId:       dc.Attributes.OwnerId,
		characterDrop: dc.Attributes.CharacterDrop,
	}
}

func attemptPickup(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c *character.Model, d *Model) {
	return func(c *character.Model, d *Model) {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		elapsed := now - int64(d.DropTime())
		if elapsed < 400 {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop has not yet met minimum time. Time elapsed %d", c.Id(), d.Id(), elapsed)
			l.Debugf("Now %d, DropTime %d, Elapsed %d.", now, d.DropTime(), elapsed)
			producers.CancelDropReservation(l, span)(d.Id(), c.Id())
			return
		}

		if !canBePickedBy(l, span)(c, d) {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop cannot be picked up by character.", c.Id(), d.Id())
			producers.CancelDropReservation(l, span)(d.Id(), c.Id())
			return
		}

		if isOwnerLockedMap(c.MapId()) && d.CharacterDrop() && d.OwnerId() != c.Id() {
			l.Debugf("Cancelling drop for character %d, drop %d, the drop is not owned by this character, in a owner locked map.", c.Id(), d.Id())
			producers.CancelDropReservation(l, span)(d.Id(), c.Id())
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
				producers.CancelDropReservation(l, span)(d.Id(), c.Id())
				// emit item unavailable.
				return
			}

			if !hasInventorySpace(l, db)(c, d) {
				l.Debugf("Cancelling drop for character %d, drop %d, the character does not have inventory space.", c.Id(), d.Id())
				producers.CancelDropReservation(l, span)(d.Id(), c.Id())
				producers.InventoryFull(l, span)(c.Id())
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
				producers.PickedUpItem(l, span)(c.Id(), d.ItemId(), d.Quantity())
			}
		}
		producers.PickedUpDrop(l, span)(c.Id(), d.Id())
	}
}

func canBePickedBy(l logrus.FieldLogger, span opentracing.Span) func(c *character.Model, d *Model) bool {
	return func(c *character.Model, d *Model) bool {
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

func pickupNX(l logrus.FieldLogger, span opentracing.Span) func(c *character.Model, d *Model) {
	return func(c *character.Model, d *Model) {
		if d.ItemId() == 4031865 {
			producers.PickedUpNx(l, span)(c.Id(), 100)
		} else {
			producers.PickedUpNx(l, span)(c.Id(), 250)
		}
	}
}

func pickupMeso(l logrus.FieldLogger, span opentracing.Span) func(c *character.Model, d *Model) {
	return func(c *character.Model, d *Model) {
		producers.AdjustMeso(l, span)(c.Id(), int32(d.Meso()))
	}
}

func consumeOnPickup() bool {
	return false
}

func needsQuestItem() bool {
	return true
}

func hasInventorySpace(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model, d *Model) bool {
	return func(c *character.Model, d *Model) bool {

		//TODO checking inventory space, and adding items should become an atomic action
		if val, ok := inventory.GetInventoryType(d.ItemId()); ok {
			i, err := inventory.GetInventoryByTypeVal(l, db)(c.Id(), val)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", c.Id())
				return false
			}

			if val == inventory.TypeValueEquip {
				return hasEquipInventorySpace(l)(c, d, i)
			} else {
				return hasItemInventorySpace(l, db)(c, d, i)
			}
		}
		return false
	}
}

func hasItemInventorySpace(l logrus.FieldLogger, db *gorm.DB) func(c *character.Model, d *Model, i *inventory.Model) bool {
	return func(c *character.Model, d *Model, i *inventory.Model) bool {

		l.Debugf("Checking inventory capacity for item %d, quantity %d for character %d.", d.ItemId(), d.Quantity(), c.Id())
		slotMax := maxInSlot()
		runningQuantity := d.Quantity()

		existingItems, err := item.GetForCharacterByInventory(l, db)(c.Id(), i.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve existing inventory %s items for character %d.", i.Type(), c.Id())
			return false
		}

		// breaks for a rechargeable item.
		usedSlots := uint32(len(existingItems))

		l.Debugf("Character %d has %d slots already occupied.", c.Id(), usedSlots)

		if len(existingItems) > 0 {
			index := 0
			for runningQuantity > 0 {
				if index < len(existingItems) {
					i := existingItems[index]
					if i.ItemId() == d.ItemId() {
						oldQuantity := i.Quantity()
						if oldQuantity < slotMax {
							newQuantity := uint32(math.Min(float64(oldQuantity+runningQuantity), float64(slotMax)))
							runningQuantity = runningQuantity - (newQuantity - oldQuantity)
						}
					}
					index++
				} else {
					break
				}
			}
		}

		newSlots := uint32(0)
		for runningQuantity > 0 {
			newQuantity := uint32(math.Min(float64(runningQuantity), float64(slotMax)))
			runningQuantity = runningQuantity - newQuantity
			newSlots += 1
		}
		l.Debugf("Character %d would need to consume %d additional slot to pick up %d %d.", c.Id(), newSlots, d.Quantity(), d.ItemId())

		if usedSlots+newSlots > i.Capacity() {
			l.Debugf("Unable to pick up item %d because character %d inventory full.", d.ItemId(), c.Id())
			return false
		}
		return true
	}
}

func hasEquipInventorySpace(l logrus.FieldLogger) func(c *character.Model, d *Model, i *inventory.Model) bool {
	return func(c *character.Model, d *Model, i *inventory.Model) bool {
		l.Debugf("Checking inventory capacity for equip %d for character %d.", d.ItemId(), c.Id())
		count := uint32(0)
		for _, equip := range i.Items() {
			if equip.Slot() >= 0 {
				count += 1
			}
		}
		if count+1 > i.Capacity() {
			l.Debugf("Unable to pick up equip %d because character %d inventory full.", d.ItemId(), c.Id())
			return false
		}
		return true
	}
}

func scriptedItem(itemId uint32) bool {
	return itemId/10000 == 243
}

func pickupEquip(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c *character.Model, d *Model) {
	return func(c *character.Model, d *Model) {
		err := equipment.GainItem(l, db, span)(c.Id(), d.ItemId(), d.EquipmentId())
		if err != nil {
			l.WithError(err).Errorf("Unable to create equipment %d that character %d picked up.", d.ItemId(), c.Id())
			return
		}
	}
}

func pickupItem(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(c *character.Model, it int8, d *Model) {
	return func(c *character.Model, it int8, d *Model) {
		err := item.GainItem(l, db, span)(c.Id(), it, d.ItemId(), d.Quantity())
		if err != nil {
			l.WithError(err).Errorf("Unable to gain item %d that character %d picked up.", d.ItemId(), c.Id())
			return
		}
	}
}

func maxInSlot() uint32 {
	return 200
}
