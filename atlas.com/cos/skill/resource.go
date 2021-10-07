package skill

import (
	"atlas-cos/json"
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{characterId}/skills", GetCharacterSkills(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{characterId}/skills/{skillId}", GetCharacterSkill(l, db)).Methods(http.MethodGet)
}

// GetCharacterSkills is a REST resource handler for retrieving the specified characters skills.
func GetCharacterSkills(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(logrus.Fields{"originator": "GetCharacterSkills", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sl, err := GetSkills(fl, db)(uint32(characterId))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get skills for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createListDataContainer(sl)

		w.WriteHeader(http.StatusOK)
		err = json.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

// GetCharacterSkill is a REST resource handler for retrieving the specified characters skill.
func GetCharacterSkill(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(logrus.Fields{"originator": "GetCharacterSkills", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		skillId, err := strconv.Atoi(mux.Vars(r)["skillId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse skillId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sl, ok := GetSkill(fl, db)(uint32(characterId), uint32(skillId))
		if !ok {
			fl.WithError(err).Errorf("Unable to get skills for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(sl)

		w.WriteHeader(http.StatusOK)
		err = json.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

func createListDataContainer(sl []*Model) attributes.CharacterSkillDataListContainer {
	var result = attributes.CharacterSkillDataListContainer{}
	result.Data = make([]attributes.CharacterSkillData, 0)
	for _, s := range sl {
		result.Data = append(result.Data, createData(s))
	}
	return result
}

func createDataContainer(s *Model) attributes.CharacterSkillDataContainer {
	var result = attributes.CharacterSkillDataContainer{}
	result.Data = createData(s)
	return result
}

func createData(model *Model) attributes.CharacterSkillData {
	return attributes.CharacterSkillData{
		Id:   strconv.Itoa(int(model.SkillId())),
		Type: "com.atlas.cos.rest.attribute.SkillAttributes",
		Attributes: attributes.CharacterSkillAttributes{
			Level:       model.Level(),
			MasterLevel: model.MasterLevel(),
			Expiration:  model.Expiration(),
		},
	}
}
