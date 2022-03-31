package location

import (
	"atlas-cos/database"
	"atlas-cos/model"
	"gorm.io/gorm"
)

// retrieveSavedLocations retrieves all saved locations for the character, or an error if one occurred.
func retrieveSavedLocations(characterId uint32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId})
	}
}

// retrieveSavedLocationsByType retrieves all saved locations for the character for the given type, or an error if one
// occurred.
func retrieveSavedLocationsByType(characterId uint32, theType string) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, &entity{CharacterId: characterId, LocationType: theType})
	}
}

// transform produces a immutable location structure
func transform(e entity) (Model, error) {
	return Model{
		id:       e.ID,
		theType:  e.LocationType,
		mapId:    e.MapId,
		portalId: e.PortalId,
	}, nil
}
