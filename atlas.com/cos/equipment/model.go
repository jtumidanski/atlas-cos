package equipment

type Model struct {
	id          uint32
	characterId uint32
	equipmentId uint32
	slot        int16
}

func (m Model) Slot() int16 {
	return m.slot
}

func (m Model) EquipmentId() uint32 {
	return m.equipmentId
}

func (m Model) Id() uint32 {
	return m.id
}
