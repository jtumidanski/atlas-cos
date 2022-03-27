package party

import (
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelProvider func() (*Model, error)

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

func ByMemberIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) ModelProvider {
	return func(memberId uint32) ModelProvider {
		return requestModelProvider(l, span)(requestByMemberId(memberId))
	}
}

func ForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (*Model, error) {
	return func(characterId uint32) (*Model, error) {
		return ByMemberIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(body requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), leaderId: att.LeaderId}
	return &m, nil
}
