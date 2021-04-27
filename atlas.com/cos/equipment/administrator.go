package equipment

import "gorm.io/gorm"

func Create(db *gorm.DB, characterId uint32, equipmentId uint32, slot int16) (*Model, error) {
	e := &entity{
		CharacterId: characterId,
		EquipmentId: equipmentId,
		Slot:        slot,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeEquipment(e), nil
}

func UpdateSlot(db *gorm.DB, equipmentId uint32, slot int16) error {
	equip, err := GetByEquipmentId(db, equipmentId)
	if err != nil {
		return err
	}

	return db.Model(&entity{Id: equip.Id()}).Select("Slot").Updates(&entity{Slot: slot}).Error
}
