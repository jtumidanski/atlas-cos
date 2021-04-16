package portal

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"log"
	"strconv"
)

type processor struct {
	l *log.Logger
}

var Processor = func(l *log.Logger) *processor {
	return &processor{l}
}

func (p *processor) GetMapPortalById(mapId uint32, portalId uint32) (*Model, error) {
	data, err := requests.MapInformation().GetPortalById(mapId, portalId)
	if err != nil {
		p.l.Printf("[ERROR] unable to get map %d portal %d.", mapId, portalId)
		return nil, err
	}
	return makePortal(data.Data()), nil
}

func (p *processor) GetMapPortals(mapId uint32) ([]*Model, error) {
	data, err := requests.MapInformation().GetPortals(mapId)
	if err != nil {
		p.l.Printf("[ERROR] unable to get map %d portals.", mapId)
		return nil, err
	}

	var portals = make([]*Model, 0)
	for _, portal := range data.DataList() {
		portals = append(portals, makePortal(&portal))
	}
	return portals, nil
}

func makePortal(data *attributes.PortalData) *Model {
	id, err := strconv.Atoi(data.Id)
	if err != nil {
		return nil
	}

	//TODO this size should be consistent, issue in POS too.
	return &Model{
		id:        uint32(id),
		theType:   data.Attributes.Type,
		x:         int16(data.Attributes.X),
		y:         int16(data.Attributes.Y),
		targetMap: data.Attributes.TargetMap,
	}
}
