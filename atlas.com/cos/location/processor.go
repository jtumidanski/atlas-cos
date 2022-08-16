package location

import (
	"atlas-cos/character"
	"atlas-cos/database"
	_map "atlas-cos/map"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetSavedLocationsByType gets the saved locations by type for the character, or returns an error if one occurred.
func GetSavedLocationsByType(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, theType string) ([]Model, error) {
	return func(characterId uint32, theType string) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(retrieveSavedLocationsByType(characterId, theType), transform)()
	}
}

// GetSavedLocations gets all the saved locations for the given character, or returns an error if one occurred.
func GetSavedLocations(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return database.ModelSliceProvider[Model, entity](db)(retrieveSavedLocations(characterId), transform)()
	}
}

// AddSavedLocation resets the saved location for of a given type for the character, returns an error if one occurred.
func AddSavedLocation(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, locationType string) error {
	return func(characterId uint32, locationType string) error {
		c, err := character.GetById(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d who is saving their location %s.", characterId, locationType)
			return err
		}

		tc := character.GetTemporalRegistry().GetById(characterId)
		if tc == nil {
			l.Errorf("Unable to locate characters %d position.", characterId)
			return errors.New("character not found")
		}

		p, err := _map.FindClosestPortal(l, span)(c.MapId(), tc.X(), tc.Y())
		if err != nil {
			l.WithError(err).Warnf("Unable to locate portal to return to. Defaulting to 0.")
		}
		portalId := p.Id()

		err = deleteByType(db, characterId, locationType)
		if err != nil {
			return err
		}
		return create(db, characterId, locationType, c.MapId(), portalId)
	}
}
