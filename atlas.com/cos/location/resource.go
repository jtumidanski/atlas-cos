package location

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/resource"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetSavedLocationsByType is a REST resource handler for retrieving the saved locations of a type for a character.
func GetSavedLocationsByType(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theType := mux.Vars(r)["type"]

		locations, err := Processor(l, db).GetSavedLocationsByType(uint32(characterId), theType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] writing GetSavedLocations response. %s", err.Error())
		}
	}
}

// GetSavedLocations is a REST resource handler for retrieving the saved locations for a character.
func GetSavedLocations(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		locations, err := Processor(l, db).GetSavedLocations(uint32(characterId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(locations)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] writing GetSavedLocations response. %s", err.Error())
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
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		li := &attributes.LocationInputDataContainer{}
		err = attributes.FromJSON(li, r.Body)
		if err != nil {
			l.Println("[ERROR] deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			err = attributes.ToJSON(&resource.GenericError{Message: err.Error()}, w)
			if err != nil {
				l.Fatalf("[ERROR] writing error message.")
			}
			return
		}

		att := li.Data.Attributes
		err = Processor(l, db).AddSavedLocation(uint32(characterId), att.Type, att.MapId, att.PortalId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
