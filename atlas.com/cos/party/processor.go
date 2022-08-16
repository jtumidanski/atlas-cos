package party

import (
	"atlas-cos/model"
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func byMemberIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) model.Provider[Model] {
	return func(memberId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestByMemberId(memberId), makeModel)
	}
}

func ForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return byMemberIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), leaderId: att.LeaderId}
	return m, nil
}
