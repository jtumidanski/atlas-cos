package information

type Model struct {
	effects []Effect
}

type Effect struct {
}

func (m Model) Effects() []Effect {
	return m.effects
}
