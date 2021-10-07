package location

import (
	"atlas-cos/json"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{characterId}/locations", HandleGetSavedLocationsByType(l, db)).Methods(http.MethodGet).Queries("type", "{type}")
	r.HandleFunc("/{characterId}/locations", HandleGetSavedLocations(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{characterId}/locations", HandleAddSavedLocation(l, db)).Methods(http.MethodPost)
}

// HandleGetSavedLocationsByType is a REST resource handler for retrieving the saved locations of a type for a character.
func HandleGetSavedLocationsByType(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(logrus.Fields{"originator": "GetSavedLocationsByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to parse characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theType := mux.Vars(r)["type"]

		locations, err := GetSavedLocationsByType(l, db)(uint32(characterId), theType)
		if err != nil {
			fl.WithError(err).Errorf("Unable to get saved locations for character %d by type %s.", characterId, theType)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = json.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

// HandleGetSavedLocations is a REST resource handler for retrieving the saved locations for a character.
func HandleGetSavedLocations(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(logrus.Fields{"originator": "GetSavedLocations", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to parse characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		locations, err := GetSavedLocations(fl, db)(uint32(characterId))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get saved locations for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = json.ToJSON(result, w)
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

// HandleAddSavedLocation is a REST resource handler for adding a saved location for a character.
func HandleAddSavedLocation(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(logrus.Fields{"originator": "AddSavedLocation", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Error retrieving characterId from request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		li := &attributes.LocationInputDataContainer{}
		err = json.FromJSON(li, r.Body)
		if err != nil {
			fl.Errorln("Deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			err = json.ToJSON(&resource.GenericError{Message: err.Error()}, w)
			if err != nil {
				fl.WithError(err).Fatalf("Writing error message.")
			}
			return
		}

		att := li.Data.Attributes
		err = AddSavedLocation(fl, db)(uint32(characterId), att.Type, att.MapId, att.PortalId)
		if err != nil {
			fl.WithError(err).Errorf("Unable to add saved location for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
