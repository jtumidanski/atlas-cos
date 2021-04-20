package character

import (
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func GetCharactersForAccountInWorld(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := strconv.Atoi(mux.Vars(r)["accountId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse accountId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId, err := strconv.Atoi(mux.Vars(r)["worldId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse worldId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cs, err := Processor(l, db).GetForAccountInWorld(uint32(accountId), byte(worldId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = &attributes.CharacterDataListContainer{}
		result.Data = make([]attributes.CharacterData, 0)
		for _, c := range cs {
			result.Data = append(result.Data, makeCharacterData(c))
		}

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}

func GetCharactersByMap(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldId, err := strconv.Atoi(mux.Vars(r)["worldId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse worldId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mapId, err := strconv.Atoi(mux.Vars(r)["mapId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse mapId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cs, err := Processor(l, db).GetForMapInWorld(byte(worldId), uint32(mapId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = &attributes.CharacterDataListContainer{}
		result.Data = make([]attributes.CharacterData, 0)
		for _, c := range cs {
			result.Data = append(result.Data, makeCharacterData(c))
		}

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}

func GetCharactersByName(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := mux.Vars(r)["name"]
		if !ok {
			l.Printf("[ERROR] unable to properly parse name from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cs, err := Processor(l, db).GetForName(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = &attributes.CharacterDataListContainer{}
		result.Data = make([]attributes.CharacterData, 0)
		for _, c := range cs {
			result.Data = append(result.Data, makeCharacterData(c))
		}

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}

func CreateCharacter(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Printf("[ERROR] unhandled request to create character.")
	}
}

func GetCharacter(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		c, err := Processor(l, db).GetById(uint32(characterId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if c == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var result = &attributes.CharacterDataContainer{}
		result.Data = makeCharacterData(c)

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}

func makeCharacterData(c *Model) attributes.CharacterData {
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

func GetCharacterDamage(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		damage := Processor(l, db).GetMaximumBaseDamage(uint32(characterId))

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
		attributes.ToJSON(result, w)
	}
}
