package equipment

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	id          uint32 `gorm:"primaryKey;autoIncrement;not null"`
	characterId uint32 `gorm:"not null"`
	equipmentId uint32 `gorm:"not null"`
	slot        int16  `gorm:"not null"`
}
