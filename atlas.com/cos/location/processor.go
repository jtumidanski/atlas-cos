package location

import (
	"gorm.io/gorm"
	"log"
)

type processor struct {
	l  *log.Logger
	db *gorm.DB
}

var Processor = func(l *log.Logger, db *gorm.DB) *processor {
	return &processor{l, db}
}

// GetSavedLocationsByType gets the saved locations by type for the character, or returns an error if one occurred.
func (p processor) GetSavedLocationsByType(characterId uint32, theType string) ([]*Model, error) {
	return getSavedLocationsByType(p.db, characterId, theType)
}

// GetSavedLocations gets all the saved locations for the given character, or returns an error if one occurred.
func (p processor) GetSavedLocations(characterId uint32) ([]*Model, error) {
	return getSavedLocations(p.db, characterId)
}

// AddSavedLocation resets the saved location for of a given type for the character, returns an error if one occurred.
func (p processor) AddSavedLocation(characterId uint32, locationType string, mapId uint32, portalId uint32) error {
	err := deleteByType(p.db, characterId, locationType)
	if err != nil {
		return nil
	}
	return create(p.db, characterId, locationType, mapId, portalId)
}
