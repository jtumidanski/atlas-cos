package portal

import (
	"atlas-cos/rest/attributes"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetMapPortalById(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalId uint32) (*Model, error) {
	return func(mapId uint32, portalId uint32) (*Model, error) {
		data, err := requestPortalById(l, span)(mapId, portalId)
		if err != nil {
			l.Errorf("Unable to get map %d portal %d.", mapId, portalId)
			return nil, err
		}
		return makePortal(data.Data()), nil
	}
}

func GetMapPortals(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]*Model, error) {
	return func(mapId uint32) ([]*Model, error) {
		data, err := requestPortals(l, span)(mapId)
		if err != nil {
			l.Errorf("Unable to get map %d portals.", mapId)
			return nil, err
		}

		var portals = make([]*Model, 0)
		for _, portal := range data.DataList() {
			portals = append(portals, makePortal(&portal))
		}
		return portals, nil
	}
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
		x:         data.Attributes.X,
		y:         data.Attributes.Y,
		targetMap: data.Attributes.TargetMapId,
	}
}
