package rest

import (
	"atlas-cos/character"
	"atlas-cos/inventory"
	"atlas-cos/location"
	"atlas-cos/seed"
	"atlas-cos/skill"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	l  *logrus.Logger
	hs *http.Server
}

func NewServer(l *logrus.Logger, db *gorm.DB) *Server {
	router := mux.NewRouter().PathPrefix("/ms/cos").Subrouter().StrictSlash(true)
	router.Use(commonHeader)

	csr := router.PathPrefix("/characters").Subrouter()
	csr.HandleFunc("", character.GetCharactersForAccountInWorld(l, db)).Methods(http.MethodGet).Queries("accountId", "{accountId}", "worldId", "{worldId}")
	csr.HandleFunc("", character.GetCharactersByMap(l, db)).Methods(http.MethodGet).Queries("worldId", "{worldId}", "mapId", "{mapId}")
	csr.HandleFunc("", character.GetCharactersByName(l, db)).Methods(http.MethodGet).Queries("name", "{name}")
	csr.HandleFunc("/{characterId}", character.GetCharacter(l, db)).Methods(http.MethodGet)
	csr.HandleFunc("/{characterId}/inventories", inventory.GetItemForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}", "slot", "{slot}")
	csr.HandleFunc("/{characterId}/inventories", inventory.GetInventoryForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}")
	csr.HandleFunc("/{characterId}/inventories/{type}/items", inventory.CreateItem(l, db)).Methods(http.MethodPost)
	csr.HandleFunc("/seeds", seed.CreateCharacterFromSeed(l, db)).Methods(http.MethodPost)
	csr.HandleFunc("/{characterId}/locations", location.HandleGetSavedLocationsByType(l, db)).Methods(http.MethodGet).Queries("type", "{type}")
	csr.HandleFunc("/{characterId}/locations", location.HandleGetSavedLocations(l, db)).Methods(http.MethodGet)
	csr.HandleFunc("/{characterId}/locations", location.HandleAddSavedLocation(l, db)).Methods(http.MethodPost)
	csr.HandleFunc("/{characterId}/damage/weapon", character.GetCharacterDamage(l, db)).Methods(http.MethodGet)
	csr.HandleFunc("/{characterId}/skills", skill.GetCharacterSkills(l, db)).Methods(http.MethodGet)
	csr.HandleFunc("/{characterId}/skills/{skillId}", skill.GetCharacterSkill(l, db)).Methods(http.MethodGet)

	w := l.Writer()
	defer w.Close()

	hs := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     log.New(w, "", 0), // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	return &Server{l, &hs}
}

func (s *Server) Run() {
	s.l.Infoln("Starting server on port 8080")
	err := s.hs.ListenAndServe()
	if err != nil {
		s.l.Errorf("Starting server: %s\n", err)
		os.Exit(1)
	}
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
