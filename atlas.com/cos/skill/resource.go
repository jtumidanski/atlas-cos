package skill

import (
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetCharacterSkills is a REST resource handler for retrieving the specified characters skills.
func GetCharacterSkills(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetCharacterSkills", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sl, err := Processor(fl, db).GetSkills(uint32(characterId))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get skills for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(sl)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			fl.WithError(err).Errorf("Writing response.")
		}
	}
}

func createDataContainer(sl []*Model) attributes.CharacterSkillDataListContainer {
	var result = attributes.CharacterSkillDataListContainer{}
	result.Data = make([]attributes.CharacterSkillData, 0)
	for _, s := range sl {
		result.Data = append(result.Data, createData(s))
	}
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
