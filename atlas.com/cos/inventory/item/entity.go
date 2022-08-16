package item

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entityInventoryItem{}, &entityItem{})
}

type entityInventoryItem struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement;not null"`
	InventoryId uint32 `gorm:"not null"`
	ItemId      uint32 `gorm:"not null"`
	Slot        int16  `gorm:"not null"`
	Type        string `gorm:"not null"`
	ReferenceId uint32
}

func (e entityInventoryItem) TableName() string {
	return "inventory_item"
}

type entityItem struct {
	ID       uint32 `gorm:"primaryKey;autoIncrement;not null"`
	Quantity uint32 `gorm:"not null"`
}

func (e entityItem) TableName() string {
	return "items"
}
