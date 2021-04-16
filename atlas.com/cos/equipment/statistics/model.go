package statistics

type Model struct {
	itemId       uint32
	weaponAttack uint16
	strength     uint16
	dexterity    uint16
	intelligence uint16
	luck         uint16
}

func (m Model) Strength() uint16 {
	return m.strength
}

func (m *Model) Dexterity() uint16 {
	return m.dexterity
}

func (m Model) Intelligence() uint16 {
	return m.intelligence
}

func (m Model) Luck() uint16 {
	return m.luck
}

func (m Model) WeaponAttack() uint16 {
	return m.weaponAttack
}

func (m Model) ItemId() uint32 {
	return m.itemId
}
