package equipment

import "gorm.io/gorm"

func Create(db *gorm.DB, characterId uint32, equipmentId uint32, slot int16) (*Model, error) {
	e := &entity{
		characterId: characterId,
		equipmentId: equipmentId,
		slot:        slot,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeEquipment(e), nil
}

func makeEquipment(e *entity) *Model {
	return &Model{
		id:          e.id,
		characterId: e.characterId,
		equipmentId: e.equipmentId,
		slot:        e.slot,
	}
}
