package rest

import (
	"atlas-cos/character"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	l  *log.Logger
	hs *http.Server
}

func NewServer(l *log.Logger, db *gorm.DB) *Server {
	router := mux.NewRouter().PathPrefix("/ms/cos").Subrouter()
	router.Use(commonHeader)

	csr := router.PathPrefix("/characters").Subrouter()
	csr.HandleFunc("", character.GetCharactersForAccountInWorld(l, db)).Methods(http.MethodGet).Queries("accountId", "{accountId}", "worldId", "{worldId}")
	csr.HandleFunc("", character.GetCharactersByMap(l, db)).Methods(http.MethodGet).Queries("worldId", "{worldId}", "mapId", "{mapId}")
	csr.HandleFunc("", character.GetCharactersByName(l, db)).Methods(http.MethodGet).Queries("name", "{name}")
	csr.HandleFunc("", character.CreateCharacter(l, db)).Methods(http.MethodPost)

	cr := csr.PathPrefix("/{characterId}").Subrouter()
	cr.HandleFunc("/{characterId}", character.GetCharacter(l, db)).Methods(http.MethodGet)

	ir := cr.PathPrefix("/inventories").Subrouter()
	ir.HandleFunc("", character.GetInventoryForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}")
	ir.HandleFunc("", character.GetInventoryForCharacter(l, db)).Methods(http.MethodGet).Queries("include", "{include}")

	seedR := cr.PathPrefix("/seeds").Subrouter()
	seedR.HandleFunc("", character.CreateCharacterFromSeed(l, db)).Methods(http.MethodPost)

	lr := cr.PathPrefix("/locations").Subrouter()
	lr.HandleFunc("", character.GetSavedLocations(l, db)).Methods(http.MethodGet).Queries("type", "{type}")
	lr.HandleFunc("", character.GetSavedLocations(l, db)).Methods(http.MethodGet)
	lr.HandleFunc("", character.AddSavedLocation(l, db)).Methods(http.MethodPost)

	dr := cr.PathPrefix("/damage").Subrouter()
	dr.HandleFunc("weapon", character.GetCharacterDamage(l, db)).Methods(http.MethodGet)

	sr := cr.PathPrefix("/skills").Subrouter()
	sr.HandleFunc("", character.GetCharacterSkills(l, db)).Methods(http.MethodGet)

	hs := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	return &Server{l, &hs}
}

func (s *Server) Run() {
	s.l.Println("[INFO] Starting server on port 8080")
	err := s.hs.ListenAndServe()
	if err != nil {
		s.l.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
