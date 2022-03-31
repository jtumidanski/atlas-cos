package character

import (
	"atlas-cos/json"
	"atlas-cos/rest"
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	GetCharactersForAccountInWorld = "get_characters_for_account_in_world"
	GetCharactersByMap             = "get_characters_by_map"
	GetCharactersByName            = "get_characters_by_name"
	GetCharacter                   = "get_character"
	GetCharacterDamage             = "get_character_damage"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("", registerGetCharactersForAccountInWorld(l, db)).Methods(http.MethodGet).Queries("accountId", "{accountId}", "worldId", "{worldId}")
	r.HandleFunc("", registerGetCharactersByMap(l, db)).Methods(http.MethodGet).Queries("worldId", "{worldId}", "mapId", "{mapId}")
	r.HandleFunc("", registerGetCharactersByName(l, db)).Methods(http.MethodGet).Queries("name", "{name}")
	r.HandleFunc("/{characterId}", registerGetCharacter(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{characterId}/damage/weapon", registerGetCharacterDamage(l, db)).Methods(http.MethodGet)
}

func registerGetCharactersForAccountInWorld(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetCharactersForAccountInWorld, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": GetCharactersForAccountInWorld, "type": "rest_handler"})
		return handleGetCharactersForAccountInWorld(fl, db)(span)
	})
}

func handleGetCharactersForAccountInWorld(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			accountId, err := strconv.Atoi(mux.Vars(r)["accountId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse accountId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			worldId, err := strconv.Atoi(mux.Vars(r)["worldId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse worldId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			cs, err := GetForAccountInWorld(l, db)(uint32(accountId), byte(worldId))
			if err != nil {
				l.WithError(err).Errorf("Unable to get characters for account %d in world %d.", accountId, worldId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result := createCharacterDataListContainer(cs)

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func registerGetCharactersByMap(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	fl := l.WithFields(logrus.Fields{"originator": GetCharactersByMap, "type": "rest_handler"})
	return rest.RetrieveSpan(GetCharactersByMap, func(span opentracing.Span) http.HandlerFunc {
		return handleGetCharactersByMap(fl, db)(span)
	})
}

func handleGetCharactersByMap(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			worldId, err := strconv.Atoi(mux.Vars(r)["worldId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse worldId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			mapId, err := strconv.Atoi(mux.Vars(r)["mapId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse mapId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			cs, err := GetForMapInWorld(l, db)(byte(worldId), uint32(mapId))
			if err != nil {
				l.WithError(err).Errorf("Unable to get characters for map %d in world %d.", mapId, worldId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result := createCharacterDataListContainer(cs)

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func registerGetCharactersByName(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	fl := l.WithFields(logrus.Fields{"originator": GetCharactersByName, "type": "rest_handler"})
	return rest.RetrieveSpan(GetCharactersByName, func(span opentracing.Span) http.HandlerFunc {
		return handleGetCharactersByName(fl, db)(span)
	})
}

func handleGetCharactersByName(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			name, ok := mux.Vars(r)["name"]
			if !ok {
				l.Errorf("Unable to properly parse name from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			cs, err := GetForName(l, db)(name)
			if err != nil {
				l.WithError(err).Errorf("Getting character %s.", name)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result := createCharacterDataListContainer(cs)

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func createCharacterDataListContainer(cs []Model) *attributes.CharacterDataListContainer {
	var result = &attributes.CharacterDataListContainer{}
	result.Data = make([]attributes.CharacterData, 0)
	for _, c := range cs {
		result.Data = append(result.Data, createCharacterData(c))
	}
	return result
}

func registerGetCharacter(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	fl := l.WithFields(logrus.Fields{"originator": GetCharacter, "type": "rest_handler"})
	return rest.RetrieveSpan(GetCharacter, func(span opentracing.Span) http.HandlerFunc {
		return handleGetCharacter(fl, db)(span)
	})
}

func handleGetCharacter(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse characterId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			c, err := GetById(l, db)(uint32(characterId))
			if err != nil {
				l.WithError(err).Errorf("Getting character %d.", characterId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var result = &attributes.CharacterDataContainer{}
			result.Data = createCharacterData(c)

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func createCharacterData(c Model) attributes.CharacterData {
	td := GetTemporalRegistry().GetById(c.Id())
	return attributes.CharacterData{
		Id:   strconv.Itoa(int(c.Id())),
		Type: "com.atlas.cos.rest.attribute.CharacterAttributes",
		Attributes: attributes.CharacterAttributes{
			AccountId:          c.AccountId(),
			WorldId:            c.WorldId(),
			Name:               c.Name(),
			Level:              c.Level(),
			Experience:         c.Experience(),
			GachaponExperience: c.GachaponExperience(),
			Strength:           c.Strength(),
			Dexterity:          c.Dexterity(),
			Intelligence:       c.Intelligence(),
			Luck:               c.Luck(),
			Hp:                 c.HP(),
			MaxHp:              c.MaxHP(),
			Mp:                 c.MP(),
			MaxMp:              c.MaxMP(),
			Meso:               c.Meso(),
			HpMpUsed:           c.HPMPUsed(),
			JobId:              c.JobId(),
			SkinColor:          c.SkinColor(),
			Gender:             c.Gender(),
			Fame:               c.Fame(),
			Hair:               c.Hair(),
			Face:               c.Face(),
			Ap:                 c.AP(),
			Sp:                 c.SPString(),
			MapId:              c.MapId(),
			SpawnPoint:         c.SpawnPoint(),
			Gm:                 c.GM(),
			X:                  td.X(),
			Y:                  td.Y(),
			Stance:             td.Stance(),
		},
	}
}

func registerGetCharacterDamage(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	fl := l.WithFields(logrus.Fields{"originator": GetCharacterDamage, "type": "rest_handler"})
	return rest.RetrieveSpan(GetCharacterDamage, func(span opentracing.Span) http.HandlerFunc {
		return handleGetCharacterDamage(fl, db)(span)
	})
}

func handleGetCharacterDamage(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
			if err != nil {
				l.WithError(err).Errorf("Unable to properly parse characterId from path.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			damage := GetMaximumBaseDamage(l, db, span)(uint32(characterId))

			result := attributes.CharacterDamageDataContainer{
				Data: attributes.CharacterDamageData{
					Id:   strconv.Itoa(characterId),
					Type: "com.atlas.cos.rest.attribute.DamageAttributes",
					Attributes: attributes.CharacterDamageAttributes{
						Type:    "WEAPON",
						Maximum: damage,
					},
				},
			}

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}
