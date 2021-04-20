package location

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/resource"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetSavedLocationsByType is a REST resource handler for retrieving the saved locations of a type for a character.
func GetSavedLocationsByType(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetSavedLocationsByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to parse characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theType := mux.Vars(r)["type"]

		locations, err := Processor(fl, db).GetSavedLocationsByType(uint32(characterId), theType)
		if err != nil {
			fl.WithError(err).Errorf("Unable to get saved locations for character %d by type %s.", characterId, theType)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

// GetSavedLocations is a REST resource handler for retrieving the saved locations for a character.
func GetSavedLocations(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetSavedLocations", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to parse characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		locations, err := Processor(fl, db).GetSavedLocations(uint32(characterId))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get saved locations for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

func createDataContainer(locations []*Model) attributes.LocationDataListContainer {
	var result attributes.LocationDataListContainer
	result.Data = make([]attributes.LocationData, 0)
	for _, loc := range locations {
		result.Data = append(result.Data, createData(loc))
	}
	return result
}

func createData(loc *Model) attributes.LocationData {
	return attributes.LocationData{
		Id:   strconv.Itoa(int(loc.Id())),
		Type: "com.atlas.cos.rest.attribute.LocationAttributes",
		Attributes: attributes.LocationAttributes{
			Type:     loc.Type(),
			MapId:    loc.MapId(),
			PortalId: loc.PortalId(),
		},
	}
}

// AddSavedLocation is a REST resource handler for adding a saved location for a character.
func AddSavedLocation(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "AddSavedLocation", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Error retrieving characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		li := &attributes.LocationInputDataContainer{}
		err = attributes.FromJSON(li, r.Body)
		if err != nil {
			fl.Errorln("Deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			err = attributes.ToJSON(&resource.GenericError{Message: err.Error()}, w)
			if err != nil {
				fl.WithError(err).Fatalf("Writing error message.")
			}
			return
		}

		att := li.Data.Attributes
		err = Processor(fl, db).AddSavedLocation(uint32(characterId), att.Type, att.MapId, att.PortalId)
		if err != nil {
			fl.WithError(err).Errorf("Unable to add saved location for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
