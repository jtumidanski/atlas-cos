package inventory

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	CharacterId   uint32 `gorm:"not null"`
	InventoryType int8   `gorm:"not null"`
	Capacity      uint32 `gorm:"capacity"`
}

func (e entity) TableName() string {
	return "inventory"
}