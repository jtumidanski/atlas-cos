package inventory

import "atlas-cos/inventory/item"

type Model struct {
	id            uint32
	inventoryType string
	capacity      uint32
	items         []item.Model
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) Type() string {
	return m.inventoryType
}

func (m Model) Capacity() uint32 {
	return m.capacity
}

func (m Model) Items() []item.Model {
	return m.items
}
