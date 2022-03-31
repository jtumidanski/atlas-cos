package seed

import (
	"atlas-cos/character"
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
	CreateCharacterFromSeed = "create_character_from_seed"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/seeds", registerCreateCharacterFromSeed(l, db)).Methods(http.MethodPost)
}

func registerCreateCharacterFromSeed(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(CreateCharacterFromSeed, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": CreateCharacterFromSeed, "type": "rest_handler"})
		return handleCreateCharacterFromSeed(fl, db)(span)
	})
}

func handleCreateCharacterFromSeed(l logrus.FieldLogger, db *gorm.DB) rest.SpanHandler {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			li := &attributes.CharacterSeedInputDataContainer{}
			err := json.FromJSON(li, r.Body)
			if err != nil {
				l.WithError(err).Errorf("Deserializing input.")
				w.WriteHeader(http.StatusBadRequest)
				err = json.ToJSON(&resource.GenericError{Message: err.Error()}, w)
				if err != nil {
					l.WithError(err).Fatalf("Writing error message.")
				}
				return
			}

			attr := li.Data.Attributes
			c, err := CreateFromSeed(l, db, span)(attr.AccountId, attr.WorldId, attr.Name, attr.JobIndex, attr.Face,
				attr.Hair, attr.HairColor, attr.Skin, attr.Gender, attr.Top, attr.Bottom, attr.Shoes, attr.Weapon)
			if err != nil {
				l.WithError(err).Errorf("Unable to create character from seed.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result := createDataContainer(c)

			w.WriteHeader(http.StatusCreated)
			err = json.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func createDataContainer(c character.Model) *attributes.CharacterDataContainer {
	var result = &attributes.CharacterDataContainer{}
	result.Data = createData(c)
	return result
}

func createData(c character.Model) attributes.CharacterData {
	td := character.GetTemporalRegistry().GetById(c.Id())
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
