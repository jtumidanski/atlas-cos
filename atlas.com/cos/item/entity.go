package item

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	id            uint32 `gorm:"primaryKey;autoIncrement;not null"`
	characterId   uint32 `gorm:"not null"`
	inventoryType byte   `gorm:"not null"`
	itemId        uint32 `gorm:"not null"`
	quantity      uint32 `gorm:"not null"`
	slot          int16  `gorm:"not null"`
}
