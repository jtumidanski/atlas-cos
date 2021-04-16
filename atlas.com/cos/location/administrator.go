package location

import "gorm.io/gorm"

func Create(db *gorm.DB, characterId uint32, locationType string, mapId uint32, portalId uint32) error {
	e := &entity{
		CharacterId:  characterId,
		LocationType: locationType,
		MapId:        mapId,
		PortalId:     portalId,
	}

	return db.Create(e).Error
}

func DeleteByType(db *gorm.DB, characterId uint32, locationType string) error {
	return db.Where("characterId = ? AND locationType =?", characterId, locationType).Delete(&entity{}).Error
}
