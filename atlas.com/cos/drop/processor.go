package drop

import (
	"atlas-cos/character"
	"atlas-cos/kafka/producers"
	"atlas-cos/party"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"context"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) AttemptPickup(characterId uint32, dropId uint32) {
	c, err := character.Processor(p.l, p.db).GetById(characterId)
	if err != nil {
		return
	}
	d, err := p.GetById(dropId)
	if err != nil {
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
	if time.Now().Sub(time.Unix(int64(d.DropTime()), 0)) < 400 {
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		return
	}

	if !p.canBePickedBy(c, d) {
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		return
	}

	if p.isOwnerLockedMap(c.MapId()) && d.CharacterDrop() && d.OwnerId() != c.Id() {
		producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
		// emit item unavailable.
		return
	}

	if d.ItemId() == 4031865 || d.ItemId() == 4031866 {
		p.pickupNX(c, d)
	} else if d.Meso() > 0 {
		p.pickupMeso(c, d)
	} else if p.consumeOnPickup(d.ItemId()) {
	} else {

		if !p.needsQuestItem(c, d) {
			producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
			// emit item unavailable.
			return
		}

		if !p.hasInventorySpace(c, d) {
			producers.CancelDropReservation(p.l, context.Background()).Emit(d.Id(), c.Id())
			// emit inventory full.
			// emit show inventory full.
			return
		}

		if p.scriptedItem(d.ItemId()) {
			// TODO handle scripted item
		}

		if val, ok := p.getInventoryType(d.ItemId()); ok {
			if val == "EQUIP" {
				p.pickupEquip(c, d)
			} else {
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
	return true
}

func (p processor) scriptedItem(itemId uint32) bool {
	return itemId/10000 == 243
}

func (p processor) getInventoryType(itemId uint32) (string, bool) {
	t := itemId / 1000000
	if t >= 1 && t <= 5 {
		switch t {
		case 1:
			return "EQUIP", true
		case 2:
			return "USE", true
		case 3:
			return "SETUP", true
		case 4:
			return "ETC", true
		case 5:
			return "CASH", true
		}
	}
	return "", false
}

func (p processor) pickupEquip(c *character.Model, d *Model) {

}

func (p processor) pickupItem(c *character.Model, d *Model, it string) {

}
