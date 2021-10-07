package statistics

import (
	"atlas-cos/rest/attributes"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func Create(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) (uint32, error) {
	return func(itemId uint32) (uint32, error) {
		ro, err := requestCreate(l, span)(itemId)
		if err != nil {
			l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
			return 0, err
		}
		eid, err := strconv.Atoi(ro.Data.Id)
		if err != nil {
			l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.")
			return 0, err
		}
		return uint32(eid), nil
	}
}

func GetEquipmentStatistics(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) (*Model, error) {
	return func(equipmentId uint32) (*Model, error) {
		resp, err := requestById(l, span)(equipmentId)
		if err != nil {
			l.WithError(err).Errorf("Retrieving equipment %d information.", equipmentId)
			return nil, err
		}
		return makeEquipment(resp.Data), nil
	}
}

func makeEquipment(resp attributes.EquipmentData) *Model {
	id, err := strconv.Atoi(resp.Id)
	if err != nil {
		return nil
	}

	attr := resp.Attributes
	return &Model{
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
	}
}
