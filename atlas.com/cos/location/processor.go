package location

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetSavedLocationsByType gets the saved locations by type for the character, or returns an error if one occurred.
func GetSavedLocationsByType(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, theType string) ([]*Model, error) {
	return func(characterId uint32, theType string) ([]*Model, error) {
		return retrieveSavedLocationsByType(db, characterId, theType)
	}
}

// GetSavedLocations gets all the saved locations for the given character, or returns an error if one occurred.
func GetSavedLocations(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		return retrieveSavedLocations(db, characterId)
	}
}

// AddSavedLocation resets the saved location for of a given type for the character, returns an error if one occurred.
func AddSavedLocation(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, locationType string, mapId uint32, portalId uint32) error {
	return func(characterId uint32, locationType string, mapId uint32, portalId uint32) error {
		err := deleteByType(db, characterId, locationType)
		if err != nil {
			return err
		}
		return create(db, characterId, locationType, mapId, portalId)
	}
}
