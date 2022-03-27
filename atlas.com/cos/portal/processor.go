package portal

import (
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

func requestModelProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelProvider {
	return func(r requests.Request[attributes]) ModelProvider {
		return func() (*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			p, err := makeModel(resp.Data())
			if err != nil {
				return nil, err
			}
			return p, nil
		}
	}
}

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelListProvider {
	return func(r requests.Request[attributes]) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(v)
				if err != nil {
					return nil, err
				}
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalId uint32) ModelProvider {
	return func(mapId uint32, portalId uint32) ModelProvider {
		return requestModelProvider(l, span)(requestById(mapId, portalId))
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalId uint32) (*Model, error) {
	return func(mapId uint32, portalId uint32) (*Model, error) {
		return ByIdModelProvider(l, span)(mapId, portalId)()
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ModelListProvider {
	return func(mapId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestInMap(mapId))
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]*Model, error) {
	return func(mapId uint32) ([]*Model, error) {
		return InMapModelProvider(l, span)(mapId)()
	}
}

func makeModel(data requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.Atoi(data.Id)
	if err != nil {
		return nil, err
	}

	attr := data.Attributes
	//TODO this size should be consistent, issue in POS too.
	return &Model{
		id:        uint32(id),
		theType:   attr.Type,
		x:         attr.X,
		y:         attr.Y,
		targetMap: attr.TargetMapId,
	}, nil
}
