package location

import "gorm.io/gorm"

type Transformer func(e *entity) *Model

// getArray queries the database for the entities which meets the provided query criteria, then applies the
// Transformer function to produce a array of results.
func getArray(db *gorm.DB, query interface{}, transformer Transformer) ([]*Model, error) {
	var result []entity
	err := db.Where(query).Find(&result).Error
	if err != nil {
		return nil, err
	}

	var skills = make([]*Model, 0)
	for _, r := range result {
		skills = append(skills, transformer(&r))
	}
	return skills, nil
}

// retrieveSavedLocations retrieves all saved locations for the character, or an error if one occurred.
func retrieveSavedLocations(db *gorm.DB, characterId uint32) ([]*Model, error) {
	return getArray(db, &entity{CharacterId: characterId}, transform)
}

// retrieveSavedLocationsByType retrieves all saved locations for the character for the given type, or an error if one
// occurred.
func retrieveSavedLocationsByType(db *gorm.DB, characterId uint32, theType string) ([]*Model, error) {
	return getArray(db, &entity{CharacterId: characterId, LocationType: theType}, transform)
}

// transform produces a immutable location structure
func transform(e *entity) *Model {
	return &Model{
		id:       e.ID,
		theType:  e.LocationType,
		mapId:    e.MapId,
		portalId: e.PortalId,
	}
}
