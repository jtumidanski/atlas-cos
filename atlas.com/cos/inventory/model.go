package inventory

type Model struct {
	id            byte
	inventoryType string
	capacity      uint32
	items         []InventoryItem
}

func (m Model) Id() byte {
	return m.id
}

func (m Model) Type() string {
	return m.inventoryType
}

func (m Model) Capacity() uint32 {
	return m.capacity
}

func (m Model) Items() []InventoryItem {
	return m.items
}

type InventoryItem struct {
	id       uint32
	itemType string
	slot     int16
}

func (i InventoryItem) Type() string {
	return i.itemType
}

func (i InventoryItem) Id() uint32 {
	return i.id
}
