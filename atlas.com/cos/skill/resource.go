package skill

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
	getCharacterSkills = "get_character_skills"
	getCharacterSkill  = "get_character_skill"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{characterId}/skills", registerGetCharacterSkills(l, db)).Methods(http.MethodGet)
	r.HandleFunc("/{characterId}/skills/{skillId}", registerGetCharacterSkill(l, db)).Methods(http.MethodGet)
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

func registerGetCharacterSkills(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(getCharacterSkills, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": getCharacterSkills, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetCharacterSkills(fl, db)(span)(characterId)
		})
	})
}

// handleGetCharacterSkills is a REST resource handler for retrieving the specified characters skills.
func handleGetCharacterSkills(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				sl, err := GetSkills(l, db)(characterId)
				if err != nil {
					l.WithError(err).Errorf("Unable to get skills for character %d.", characterId)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				result := createListDataContainer(sl)

				w.WriteHeader(http.StatusOK)
				err = json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Writing response.")
				}
			}
		}
	}
}

func registerGetCharacterSkill(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(getCharacterSkill, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": getCharacterSkill, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return parseSkillId(fl, func(skillId uint32) http.HandlerFunc {
				return handleGetCharacterSkill(fl, db)(span)(characterId)(skillId)
			})
		})
	})
}

type skillIdHandler func(skillId uint32) http.HandlerFunc

func parseSkillId(l logrus.FieldLogger, next skillIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skillId, err := strconv.Atoi(mux.Vars(r)["skillId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse skillId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(skillId))(w, r)
	}
}

// handleGetCharacterSkill is a REST resource handler for retrieving the specified character's skill.
func handleGetCharacterSkill(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) func(skillId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) func(skillId uint32) http.HandlerFunc {
		return func(characterId uint32) func(skillId uint32) http.HandlerFunc {
			return func(skillId uint32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					sl, ok := GetSkill(l, db)(characterId, skillId)
					if !ok {
						l.Errorf("Unable to get skills for character %d.", characterId)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					result := createDataContainer(sl)

					w.WriteHeader(http.StatusOK)
					err := json.ToJSON(result, w)
					if err != nil {
						l.WithError(err).Errorf("Writing response.")
					}
				}
			}
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
