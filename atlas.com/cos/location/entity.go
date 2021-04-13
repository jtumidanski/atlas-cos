package location

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	ID           uint32 `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId  uint32 `gorm:"not null"`
	LocationType string `gorm:"not null"`
	MapId        uint32 `gorm:"not null"`
	PortalId     uint32 `gorm:"not null"`
}
