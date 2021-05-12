package statistics

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetEquipmentStatistics(l logrus.FieldLogger) func(equipmentId uint32) (*Model, error) {
	return func(equipmentId uint32) (*Model, error) {
		resp, err := requests.EquipmentRegistry().GetById(equipmentId)
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
