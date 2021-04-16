package item

type Model struct {
	id            uint32
	characterId   uint32
	inventoryType byte
	itemId        uint32
	quantity      uint32
	slot          int16
}

func (m Model) Quantity() uint32 {
	return m.quantity
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) Slot() int16 {
	return m.slot
}

func (m Model) InventoryType() byte {
	return m.inventoryType
}
