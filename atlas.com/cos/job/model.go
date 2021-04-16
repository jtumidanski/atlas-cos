package job

type Model struct {
	id uint16
}

func (m Model) Id() uint16 {
	return m.id
}
