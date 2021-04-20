package location

import "gorm.io/gorm"

// create creates a saved location for a character with the provided attributes, and returns an error if one occurred.
func create(db *gorm.DB, characterId uint32, locationType string, mapId uint32, portalId uint32) error {
	e := &entity{
		CharacterId:  characterId,
		LocationType: locationType,
		MapId:        mapId,
		PortalId:     portalId,
	}

	return db.Create(e).Error
}

// deleteByType deletes all saved locations for a character with the given type.
func deleteByType(db *gorm.DB, characterId uint32, locationType string) error {
	return db.Where(&entity{CharacterId: characterId, LocationType: locationType}).Delete(&entity{}).Error
}
