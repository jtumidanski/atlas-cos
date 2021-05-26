package rest

import (
	"atlas-cos/character"
	"atlas-cos/inventory"
	"atlas-cos/location"
	"atlas-cos/seed"
	"atlas-cos/skill"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, db *gorm.DB, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes(db))
}

func ProduceRoutes(db *gorm.DB) func(l logrus.FieldLogger) http.Handler {
	return func(l logrus.FieldLogger) http.Handler {
		router := mux.NewRouter().PathPrefix("/ms/cos").Subrouter().StrictSlash(true)
		router.Use(CommonHeader)

		csr := router.PathPrefix("/characters").Subrouter()
		csr.HandleFunc("", character.GetCharactersForAccountInWorld(l, db)).Methods(http.MethodGet).Queries("accountId", "{accountId}", "worldId", "{worldId}")
		csr.HandleFunc("", character.GetCharactersByMap(l, db)).Methods(http.MethodGet).Queries("worldId", "{worldId}", "mapId", "{mapId}")
		csr.HandleFunc("", character.GetCharactersByName(l, db)).Methods(http.MethodGet).Queries("name", "{name}")
		csr.HandleFunc("/{characterId}", character.GetCharacter(l, db)).Methods(http.MethodGet)
		csr.HandleFunc("/{characterId}/inventories", inventory.GetItemForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}", "slot", "{slot}")
		csr.HandleFunc("/{characterId}/inventories", inventory.GetItemsForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}", "itemId", "{itemId}")
		csr.HandleFunc("/{characterId}/inventories", inventory.GetInventoryForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}")
		csr.HandleFunc("/{characterId}/inventories/{type}/items", inventory.CreateItem(l, db)).Methods(http.MethodPost)
		csr.HandleFunc("/{characterId}/items", inventory.GetItemsForCharacter(l, db)).Methods(http.MethodGet).Queries("itemId", "{itemId}")
		csr.HandleFunc("/seeds", seed.CreateCharacterFromSeed(l, db)).Methods(http.MethodPost)
		csr.HandleFunc("/{characterId}/locations", location.HandleGetSavedLocationsByType(l, db)).Methods(http.MethodGet).Queries("type", "{type}")
		csr.HandleFunc("/{characterId}/locations", location.HandleGetSavedLocations(l, db)).Methods(http.MethodGet)
		csr.HandleFunc("/{characterId}/locations", location.HandleAddSavedLocation(l, db)).Methods(http.MethodPost)
		csr.HandleFunc("/{characterId}/damage/weapon", character.GetCharacterDamage(l, db)).Methods(http.MethodGet)
		csr.HandleFunc("/{characterId}/skills", skill.GetCharacterSkills(l, db)).Methods(http.MethodGet)
		csr.HandleFunc("/{characterId}/skills/{skillId}", skill.GetCharacterSkill(l, db)).Methods(http.MethodGet)

		return router
	}
}
