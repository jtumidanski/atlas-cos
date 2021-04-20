package skill

import (
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetCharacterSkills is a REST resource handler for retrieving the specified characters skills.
func GetCharacterSkills(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sl, err := Processor(l, db).GetSkills(uint32(characterId))
		if err != nil {
			l.Printf("[ERROR] %s.", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = attributes.CharacterSkillDataListContainer{}
		result.Data = make([]attributes.CharacterSkillData, 0)
		for _, s := range sl {
			result.Data = append(result.Data, createCharacterSkillData(s))
		}

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] writing GetCharacterSkills response. %s", err.Error())
		}
	}
}

func createCharacterSkillData(model *Model) attributes.CharacterSkillData {
	return attributes.CharacterSkillData{
		Id:   strconv.Itoa(int(model.Id())),
		Type: "com.atlas.cos.rest.attribute.SkillAttributes",
		Attributes: attributes.CharacterSkillAttributes{
			Level:       model.Level(),
			MasterLevel: model.MasterLevel(),
			Expiration:  model.Expiration(),
		},
	}
}
