package portal

type Model struct {
	id        uint32
	theType   uint32
	x         int16
	y         int16
	targetMap uint32
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) X() int16 {
	return m.x
}

func (m Model) Y() int16 {
	return m.y
}

func (m Model) IsSpawnPoint() bool {
	return (m.theType == 0 || m.theType == 1) && m.targetMap == 999999999
}
