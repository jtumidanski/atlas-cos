package item

type Model interface {
	Id() uint32
	InventoryId() uint32
	ItemId() uint32
	Slot() int16
}

type EquipmentModel struct {
	id          uint32
	inventoryId uint32
	itemId      uint32
	slot        int16
	equipmentId uint32
}

func (m EquipmentModel) Id() uint32 {
	return m.id
}

func (m EquipmentModel) InventoryId() uint32 {
	return m.inventoryId
}

func (m EquipmentModel) ItemId() uint32 {
	return m.itemId
}

func (m EquipmentModel) EquipmentId() uint32 {
	return m.equipmentId
}

func (m EquipmentModel) Slot() int16 {
	return m.slot
}

type ItemModel struct {
	id          uint32
	inventoryId uint32
	itemId      uint32
	slot        int16
	quantity    uint32
	referenceId uint32
}

func (m ItemModel) Id() uint32 {
	return m.id
}

func (m ItemModel) InventoryId() uint32 {
	return m.inventoryId
}

func (m ItemModel) ItemId() uint32 {
	return m.itemId
}

func (m ItemModel) Slot() int16 {
	return m.slot
}

func (m ItemModel) Quantity() uint32 {
	return m.quantity
}

func (m ItemModel) ReferenceId() uint32 {
	return m.referenceId
}
