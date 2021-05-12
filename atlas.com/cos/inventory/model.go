package inventory

type Model struct {
	id            int8
	inventoryType string
	capacity      uint32
	items         []Item
}

func (m Model) Id() int8 {
	return m.id
}

func (m Model) Type() string {
	return m.inventoryType
}

func (m Model) Capacity() uint32 {
	return m.capacity
}

func (m Model) Items() []Item {
	return m.items
}

type Item struct {
	id       uint32
	itemType string
	slot     int16
}

func (i Item) Type() string {
	return i.itemType
}

func (i Item) Id() uint32 {
	return i.id
}

func (i Item) Slot() int16 {
	return i.slot
}
