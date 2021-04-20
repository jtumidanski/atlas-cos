package seed

import (
	"atlas-cos/character"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/resource"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func CreateCharacterFromSeed(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		li := &attributes.CharacterSeedInputDataContainer{}
		err := attributes.FromJSON(li, r.Body)
		if err != nil {
			l.Println("[ERROR] deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			err = attributes.ToJSON(&resource.GenericError{Message: err.Error()}, w)
			if err != nil {
				l.Fatalf("[ERROR] writing error message.")
			}
			return
		}

		attr := li.Data.Attributes
		c, err := Processor(l, db).CreateFromSeed(attr.AccountId, attr.WorldId, attr.Name, attr.JobIndex, attr.Face,
			attr.Hair, attr.HairColor, attr.Skin, attr.Gender, attr.Top, attr.Bottom, attr.Shoes, attr.Weapon)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createDataContainer(c)

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] writing CreateCharacterFromSeed response. %s", err.Error())
		}
	}
}

func createDataContainer(c *character.Model) *attributes.CharacterDataContainer {
	var result = &attributes.CharacterDataContainer{}
	result.Data = createData(c)
	return result
}

func createData(c *character.Model) attributes.CharacterData {
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
