package statistics

import (
	"atlas-cos/model"
	"atlas-cos/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func Create(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) (uint32, error) {
	return func(itemId uint32) (uint32, error) {
		ro, _, err := requestCreate(itemId)(l, span)
		if err != nil {
			l.WithError(err).Errorf("Generating equipment item %d, they were not awarded this item. Check request in ESO service.", itemId)
			return 0, err
		}
		eid, err := strconv.Atoi(ro.Data().Id)
		if err != nil {
			l.WithError(err).Errorf("Generating equipment item %d, they were not awarded this item. Invalid ID from ESO service.", itemId)
			return 0, err
		}
		return uint32(eid), nil
	}
}

func byEquipmentIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) model.Provider[Model] {
	return func(equipmentId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(equipmentId), makeEquipment)
	}
}

func GetEquipmentStatistics(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) (Model, error) {
	return func(equipmentId uint32) (Model, error) {
		return byEquipmentIdModelProvider(l, span)(equipmentId)()
	}
}

func makeEquipment(resp requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.Atoi(resp.Id)
	if err != nil {
		return Model{}, err
	}

	attr := resp.Attributes
	return Model{
		id:            uint32(id),
		itemId:        attr.ItemId,
		strength:      attr.Strength,
		dexterity:     attr.Dexterity,
		intelligence:  attr.Intelligence,
		luck:          attr.Luck,
		hp:            attr.HP,
		mp:            attr.MP,
		weaponAttack:  attr.WeaponAttack,
		magicAttack:   attr.MagicAttack,
		weaponDefense: attr.WeaponDefense,
		magicDefense:  attr.MagicDefense,
		accuracy:      attr.Accuracy,
		avoidability:  attr.Avoidability,
		hands:         attr.Hands,
		speed:         attr.Speed,
		jump:          attr.Jump,
		slots:         attr.Slots,
		cash:          attr.Cash,
	}, nil
}
