package skill

type Model struct {
	id          uint32
	level       uint32
	masterLevel uint32
	expiration  int64
}

func (m Model) Level() uint32 {
	return m.level
}

func (m Model) MasterLevel() uint32 {
	return m.masterLevel
}

func (m Model) Expiration() int64 {
	return m.expiration
}

func (m Model) Id() uint32 {
	return m.id
}
