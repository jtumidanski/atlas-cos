package equipment

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId uint32 `gorm:"not null"`
	EquipmentId uint32 `gorm:"not null"`
	Slot        int16  `gorm:"not null"`
}
