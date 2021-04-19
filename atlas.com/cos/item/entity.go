package item

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	Id            uint32 `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId   uint32 `gorm:"not null"`
	InventoryType byte   `gorm:"not null"`
	ItemId        uint32 `gorm:"not null"`
	Quantity      uint32 `gorm:"not null"`
	Slot          int16  `gorm:"not null"`
}

func (e entity) TableName() string {
	return "items"
}
