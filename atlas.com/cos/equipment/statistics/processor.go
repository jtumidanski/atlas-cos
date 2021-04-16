package statistics

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

func (p processor) GetEquipmentStatistics(equipmentId uint32) (*Model, error) {
	resp, err := requests.EquipmentRegistry().GetById(equipmentId)
	if err != nil {
		p.l.Printf("[ERROR] retrieving equipment %d information.", equipmentId)
		return nil, err
	}
	return makeEquipment(resp.Data), nil
}

func makeEquipment(resp attributes.EquipmentData) *Model {
	return &Model{
		itemId:       resp.Attributes.ItemId,
		weaponAttack: resp.Attributes.WeaponAttack,
		strength:     resp.Attributes.Strength,
		dexterity:    resp.Attributes.Dexterity,
		intelligence: resp.Attributes.Intelligence,
		luck:         resp.Attributes.Luck,
	}
}
