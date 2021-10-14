package location

import (
	"atlas-cos/json"
	"atlas-cos/rest"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/resource"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	getSavedLocationsByType = "get_saved_locations_by_type"
	getSavedLocations       = "get_saved_locations"
	addSavedLocation        = "add_saved_location"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{characterId}/locations", registerGetSavedLocationsByType(l, db)).Methods(http.MethodGet).Queries("type", "{type}")
	r.HandleFunc("/{characterId}/locations", registerGetSavedLocations(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{characterId}/locations", registerAddSavedLocation(l, db)).Methods(http.MethodPost)
}

type characterIdHandler func(characterId uint32) http.HandlerFunc

func parseCharacterId(l logrus.FieldLogger, next characterIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(characterId))(w, r)
	}
}

func registerGetSavedLocationsByType(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(getSavedLocationsByType, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": getSavedLocationsByType, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetSavedLocationsByType(fl, db)(span)(characterId)
		})
	})
}

// handleGetSavedLocationsByType is a REST resource handler for retrieving the saved locations of a type for a character.
func handleGetSavedLocationsByType(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				theType := mux.Vars(r)["type"]
				locations, err := GetSavedLocationsByType(l, db)(characterId, theType)
				if err != nil {
					l.WithError(err).Errorf("Unable to get saved locations for character %d by type %s.", characterId, theType)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				result := createDataContainer(locations)

				w.WriteHeader(http.StatusOK)
				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Writing response.")
				}
			}
		}
	}
}

func registerGetSavedLocations(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(getSavedLocations, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": getSavedLocations, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetSavedLocations(fl, db)(span)(characterId)
		})
	})
}

// handleGetSavedLocations is a REST resource handler for retrieving the saved locations for a character.
func handleGetSavedLocations(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				locations, err := GetSavedLocations(l, db)(characterId)
				if err != nil {
					l.WithError(err).Errorf("Unable to get saved locations for character %d.", characterId)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				result := createDataContainer(locations)

				w.WriteHeader(http.StatusOK)
				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Writing response.")
				}
			}
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

func registerAddSavedLocation(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(addSavedLocation, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": addSavedLocation, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleAddSavedLocation(fl, db)(span)(characterId)
		})
	})
}

// handleAddSavedLocation is a REST resource handler for adding a saved location for a character.
func handleAddSavedLocation(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				li := &attributes.LocationInputDataContainer{}
				err := json.FromJSON(li, r.Body)
				if err != nil {
					l.Errorln("Deserializing input", err)
					w.WriteHeader(http.StatusBadRequest)
					err = json.ToJSON(&resource.GenericError{Message: err.Error()}, w)
					if err != nil {
						l.WithError(err).Fatalf("Writing error message.")
					}
					return
				}

				att := li.Data.Attributes
				err = AddSavedLocation(l, db)(characterId, att.Type, att.MapId, att.PortalId)
				if err != nil {
					l.WithError(err).Errorf("Unable to add saved location for character %d.", characterId)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
	}
}
