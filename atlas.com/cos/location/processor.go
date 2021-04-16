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

func (p processor) GetSavedLocationsByType(characterId uint32, theType string) ([]*Model, error) {
	return GetSavedLocationsByType(p.db, characterId, theType)
}

func (p processor) GetSavedLocations(characterId uint32) ([]*Model, error) {
	return GetSavedLocations(p.db, characterId)
}

func (p processor) AddSavedLocation(characterId uint32, locationType string, mapId uint32, portalId uint32) error {
	err := DeleteByType(p.db, characterId, locationType)
	if err != nil {
		return nil
	}
	return Create(p.db, characterId, locationType, mapId, portalId)
}
