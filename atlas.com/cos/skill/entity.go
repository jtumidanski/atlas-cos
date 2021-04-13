package skill

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	_ = db.AutoMigrate(&entity{})
}

type entity struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement;not null"`
	SkillId     uint32 `gorm:"not null"`
	CharacterId uint32 `gorm:"not null"`
	SkillLevel  uint32 `gorm:"not null"`
	MasterLevel uint32 `gorm:"not null"`
	Expiration  uint64 `gorm:"not null;default=-1"`
}
