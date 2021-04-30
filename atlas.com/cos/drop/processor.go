package drop

import (
	"atlas-cos/character"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/kafka/producers"
	"atlas-cos/party"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type processor struct {
	l  log.FieldLogger
	db *gorm.DB
}

var Processor = func(l log.FieldLogger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) AttemptPickup(characterId uint32, dropId uint32) {
	c, err := character.Processor(p.l, p.db).GetById(characterId)
	if err != nil {
		p.l.WithError(err).Errorf("Attempting to pick up %d for character %d.", dropId, characterId)
		return
	}
	d, err := p.GetById(dropId)
	if err != nil {
		p.l.WithError(err).Errorf("Attempting to pick up %d for character %d.", dropId, characterId)
		return
	}
	p.attemptPickup(c, d)
}

func (p processor) GetById(dropId uint32) (*Model, error) {
	dc, err := requests.DropRegistry().GetDropById(dropId)
	if err != nil {
		return nil, err
	}
	d := p.makeDrop(dc.Data)
	return d, nil
}

func (p processor) makeDrop(dc attributes.DropData) *Model {
	id, err := strconv.Atoi(dc.Id)
	if err != nil {
		return nil
	}
	return &Model{
		id:            uint32(id),
		itemId:        dc.Attributes.ItemId,
		quantity:      dc.Attributes.Quantity,
		meso:          dc.Attributes.Meso,
		dropTime:      dc.Attributes.DropTime,
		dropType:      dc.Attributes.DropType,
		ownerId:       dc.Attributes.OwnerId,
		characterDrop: dc.Attributes.CharacterDrop,
	}
}

func (p processor) attemptPickup(c *character.Model, d *Model) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	elapsed := now - int64(d.DropTime())
	if elapsed < 400 {
		p.l.Debugf("Cancelling drop for character %d, drop %d, the drop has not yet met minimum time. Time elapsed %d", c.Id(), d.Id(), elapsed)
		p.l.Debugf("Now %d, DropTime %d, Elapsed %d.", now, d.DropTime(), elapsed)
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		return
	}

	if !p.canBePickedBy(c, d) {
		p.l.Debugf("Cancelling drop for character %d, drop %d, the drop cannot be picked up by character.", c.Id(), d.Id())
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		return
	}

	if p.isOwnerLockedMap(c.MapId()) && d.CharacterDrop() && d.OwnerId() != c.Id() {
		p.l.Debugf("Cancelling drop for character %d, drop %d, the drop is not owned by this character, in a owner locked map.", c.Id(), d.Id())
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		// emit item unavailable.
		return
	}

	if d.ItemId() == 4031865 || d.ItemId() == 4031866 {
		p.l.Debugf("Picking up NX item %d for character %d.", d.ItemId(), c.Id())
		p.pickupNX(c, d)
	} else if d.Meso() > 0 {
		p.l.Debugf("Picking up meso drop for character %d.", c.Id())
		p.pickupMeso(c, d)
	} else if p.consumeOnPickup(d.ItemId()) {
	} else {

		if !p.needsQuestItem(c, d) {
			p.l.Debugf("Cancelling drop for character %d, drop %d, the character does not need this quest item.", c.Id(), d.Id())
			producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
			// emit item unavailable.
			return
		}

		if !p.hasInventorySpace(c, d) {
			p.l.Debugf("Cancelling drop for character %d, drop %d, the character does not have inventory space.", c.Id(), d.Id())
			producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
			// emit inventory full.
			// emit show inventory full.
			return
		}

		if p.scriptedItem(d.ItemId()) {
			// TODO handle scripted item
		}

		if val, ok := inventory.GetInventoryType(d.ItemId()); ok {
			if val == inventory.TypeValueEquip {
				p.l.Debugf("Picking up equip item %d for character %d.", d.ItemId(), c.Id())
				p.pickupEquip(c, d)
			} else {
				p.l.Debugf("Picking up item %d for character %d.", d.ItemId(), c.Id())
				p.pickupItem(c, d, val)
			}
			// TODO update ariant score if 4031868
			producers.PickedUpItem(p.l, context.Background()).Emit(c.Id(), d.ItemId(), d.Quantity())
		}
	}
	producers.PickedUpDrop(p.l, context.Background()).Emit(c.Id(), d.Id())
}

func (p processor) canBePickedBy(c *character.Model, d *Model) bool {
	if d.OwnerId() <= 0 || d.FFADrop() {
		return true
	}

	ownerParty, err := party.Processor(p.l, p.db).PartyForCharacter(d.OwnerId())
	if err != nil {
		if c.Id() == d.OwnerId() {
			return true
		}
	} else {
		characterParty, err := party.Processor(p.l, p.db).PartyForCharacter(c.Id())
		if err == nil && ownerParty.Id() == characterParty.Id() {
			return true
		} else if c.Id() == d.OwnerId() {
			return true
		}
	}
	return d.ExpiredOwnershipTime()
}

func (p processor) isOwnerLockedMap(mapId uint32) bool {
	return mapId > 209000000 && mapId < 209000016 || mapId >= 990000500 && mapId <= 990000502
}

func (p processor) pickupNX(c *character.Model, d *Model) {
	if d.ItemId() == 4031865 {
		producers.PickedUpNx(p.l, context.Background()).Emit(c.Id(), 100)
	} else {
		producers.PickedUpNx(p.l, context.Background()).Emit(c.Id(), 250)
	}
}

func (p processor) pickupMeso(c *character.Model, d *Model) {
	producers.AdjustMeso(p.l, context.Background()).Emit(c.Id(), d.Meso())
}

func (p processor) consumeOnPickup(itemId uint32) bool {
	return false
}

func (p processor) needsQuestItem(c *character.Model, d *Model) bool {
	return true
}

func (p processor) hasInventorySpace(c *character.Model, d *Model) bool {
	//TODO checking inventory space, and adding items should become an atomic action
	if val, ok := inventory.GetInventoryType(d.ItemId()); ok {
		i, err := inventory.Processor(p.l, p.db).GetInventoryByTypeVal(c.Id(), val)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", c.Id())
			return false
		}

		if val == inventory.TypeValueEquip {
			return p.hasEquipInventorySpace(c, d, i)
		} else {
			return p.hasItemInventorySpace(c, d, i)
		}
	}
	return false
}

func (p processor) hasItemInventorySpace(c *character.Model, d *Model, i *inventory.Model) bool {
	p.l.Debugf("Checking inventory capacity for item %d, quantity %d for character %d.", d.ItemId(), d.Quantity(), c.Id())
	slotMax := p.maxInSlot(c, d)
	runningQuantity := d.Quantity()

	existingItems, err := item.Processor(p.l, p.db).GetForCharacterByInventory(c.Id(), i.Id())
	if err != nil {
		p.l.WithError(err).Errorf("Unable to retrieve existing inventory %d items for character %d.", i.Type(), c.Id())
		return false
	}

	// breaks for a rechargeable item.
	usedSlots := uint32(len(existingItems))

	p.l.Debugf("Character %d has %d slots already occupied.", c.Id(), usedSlots)

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
	p.l.Debugf("Character %d would need to consume %d additional slot to pick up %d %d.", c.Id(), newSlots, d.Quantity(), d.ItemId())

	if usedSlots+newSlots > i.Capacity() {
		p.l.Debugf("Unable to pick up item %d because character %d inventory full.", d.ItemId(), c.Id())
		return false
	}
	return true
}

func (p processor) hasEquipInventorySpace(c *character.Model, d *Model, i *inventory.Model) bool {
	p.l.Debugf("Checking inventory capacity for equip %d for character %d.", d.ItemId(), c.Id())
	count := uint32(0)
	for _, equip := range i.Items() {
		if equip.Slot() >= 0 {
			count += 1
		}
	}
	if count+1 > i.Capacity() {
		p.l.Debugf("Unable to pick up equip %d because character %d inventory full.", d.ItemId(), c.Id())
		return false
	}
	return true
}

func (p processor) scriptedItem(itemId uint32) bool {
	return itemId/10000 == 243
}

func (p processor) pickupEquip(c *character.Model, d *Model) {
	e, err := equipment.Processor(p.l, p.db).CreateForCharacter(c.Id(), d.ItemId(), false)
	if err != nil {
		p.l.WithError(err).Errorf("Unable to create equipment %d that character %d picked up.", d.ItemId(), c.Id())
		return
	}
	producers.InventoryModificationReservation(p.l, context.Background()).
		Emit(c.Id(), true, 0, d.ItemId(), 1, d.Quantity(), e.Slot(), 0)
}

func (p processor) pickupItem(c *character.Model, d *Model, it int8) {
	slotMax := p.maxInSlot(c, d)
	runningQuantity := d.Quantity()

	existingItems := item.Processor(p.l, p.db).GetItemsForCharacter(c.Id(), it, d.ItemId())
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
					err := item.Processor(p.l, p.db).UpdateItemQuantity(i.Id(), newQuantity)
					if err != nil {
						p.l.WithError(err).Errorf("Updating the quantity of item %d to value %d.", i.Id(), newQuantity)
					} else {
						producers.InventoryModificationReservation(p.l, context.Background()).
							Emit(c.Id(), true, 1, d.ItemId(), i.InventoryType(), newQuantity, i.Slot(), 0)
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
		i, err := item.Processor(p.l, p.db).CreateItemForCharacter(c.Id(), it, d.ItemId(), newQuantity)
		if err != nil {
			p.l.WithError(err).Errorf("Unable to create item %d that character %d picked up.", d.ItemId(), c.Id())
			return
		}
		producers.InventoryModificationReservation(p.l, context.Background()).
			Emit(c.Id(), true, 0, d.ItemId(), it, d.Quantity(), i.Slot(), 0)
	}
}

func (p processor) maxInSlot(c *character.Model, d *Model) uint32 {
	return 200
}
