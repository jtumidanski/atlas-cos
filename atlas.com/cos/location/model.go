package location

type Model struct {
	id       uint32
	theType  string
	mapId    uint32
	portalId uint32
}

func (m Model) Type() string {
	return m.theType
}

func (m Model) MapId() uint32 {
	return m.mapId
}

func (m Model) PortalId() uint32 {
	return m.portalId
}

func (m Model) Id() uint32 {
	return m.id
}
