package drop

import "time"

type Model struct {
	id            uint32
	itemId        uint32
	equipmentId   uint32
	quantity      uint32
	meso          uint32
	dropTime      uint64
	dropType      byte
	ownerId       uint32
	characterDrop bool
}

func (m Model) DropTime() uint64 {
	return m.dropTime
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) CharacterDrop() bool {
	return m.characterDrop
}

func (m Model) OwnerId() uint32 {
	return m.ownerId
}

func (m Model) FFADrop() bool {
	return m.dropType == 2 || m.dropType == 3 || m.ExpiredOwnershipTime()
}

func (m Model) ExpiredOwnershipTime() bool {
	return time.Now().Sub(time.Unix(int64(m.dropTime), 0)) >= time.Millisecond*time.Duration(15)
}

func (m Model) ItemId() uint32 {
	return m.itemId
}

func (m Model) Meso() uint32 {
	return m.meso
}

func (m Model) Quantity() uint32 {
	return m.quantity
}

func (m Model) EquipmentId() uint32 {
	return m.equipmentId
}
