package portal

import (
	"atlas-cos/model"
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalId uint32) model.Provider[Model] {
	return func(mapId uint32, portalId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(mapId, portalId), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalId uint32) (Model, error) {
	return func(mapId uint32, portalId uint32) (Model, error) {
		return ByIdModelProvider(l, span)(mapId, portalId)()
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.SliceProvider[Model] {
	return func(mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(mapId), makeModel)
	}
}

func makeModel(data requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.Atoi(data.Id)
	if err != nil {
		return Model{}, err
	}

	attr := data.Attributes
	//TODO this size should be consistent, issue in POS too.
	return Model{
		id:        uint32(id),
		theType:   attr.Type,
		x:         attr.X,
		y:         attr.Y,
		targetMap: attr.TargetMapId,
	}, nil
}
