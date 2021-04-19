package location

import "gorm.io/gorm"

func GetSavedLocations(db *gorm.DB, characterId uint32) ([]*Model, error) {
	var results []entity
	err := db.Where(&entity{CharacterId: characterId}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var locations = make([]*Model, 0)
	for _, e := range results {
		locations = append(locations, makeLocation(&e))
	}
	return locations, nil
}

func makeLocation(e *entity) *Model {
	return &Model{
		id:       e.ID,
		theType:  e.LocationType,
		mapId:    e.MapId,
		portalId: e.PortalId,
	}
}

func GetSavedLocationsByType(db *gorm.DB, characterId uint32, theType string) ([]*Model, error) {
	var results []entity
	err := db.Where(&entity{CharacterId: characterId, LocationType: theType}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var locations = make([]*Model, 0)
	for _, e := range results {
		locations = append(locations, makeLocation(&e))
	}
	return locations, nil
}
