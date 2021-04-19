package character

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/location"
	"atlas-cos/rest/attributes"
	"atlas-cos/skill"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func GetCharactersForAccountInWorld(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
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

func GetCharactersByMap(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
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

func GetCharactersByName(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
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

func CreateCharacter(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Printf("[ERROR] unhandled request to create character.")
	}
}

func GetCharacter(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
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

func GetInventoryForCharacterByType(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		include := mux.Vars(r)["include"]
		inventoryType := mux.Vars(r)["type"]
		if inventoryType == "" {
			l.Printf("[ERROR] unable to retrieve requested inventory type.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inv, err := inventory.Processor(l, db).GetInventoryByType(uint32(characterId), inventoryType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = &attributes.InventoryDataContainer{
			Data: attributes.InventoryData{
				Id:   strconv.Itoa(int(inv.Id())),
				Type: "com.atlas.cos.rest.attribute.InventoryAttributes",
				Attributes: attributes.InventoryAttributes{
					Type:     inv.Type(),
					Capacity: inv.Capacity(),
				},
			},
		}

		for _, i := range inv.Items() {
			var inventoryItem = attributes.Relationship{
				Id:   strconv.Itoa(int(i.Id())),
				Type: getInventoryItemType(i.Type()),
			}

			result.Data.Relationships.InventoryItems = append(result.Data.Relationships.InventoryItems, inventoryItem)
		}

		if strings.Contains(include, "inventoryItems") {
			inv, err := inventory.Processor(l, db).GetInventoryByType(uint32(characterId), inventoryType)
			if err != nil {
			} else {
				for _, inventoryItem := range inv.Items() {
					if inventoryItem.Type() == inventory.ItemTypeEquip {
						e, err := equipment.Processor(l, db).GetEquipmentById(inventoryItem.Id())
						if err != nil {
						} else {
							result.Included = append(result.Included, attributes.InventoryEquipmentData{
								Id:   strconv.Itoa(int(e.Id())),
								Type: "com.atlas.cos.rest.attribute.EquipmentAttributes",
								Attributes: attributes.InventoryEquipmentAttributes{
									EquipmentId: e.EquipmentId(),
									Slot:        e.Slot(),
								},
								Relationships: attributes.Relationships{
									EquipmentStatistics: []attributes.Relationship{{
										Id:   strconv.Itoa(int(e.EquipmentId())),
										Type: "com.atlas.cos.rest.attribute.EquipmentStatisticsAttributes",
									}},
								},
							})
						}
					} else {
						i, err := item.Processor(l, db).GetItemById(inventoryItem.Id())
						if err != nil {
						} else {
							result.Included = append(result.Included, attributes.InventoryItemData{
								Id:   strconv.Itoa(int(i.Id())),
								Type: "com.atlas.cos.rest.attribute.ItemAttributes",
								Attributes: attributes.InventoryItemAttributes{
									ItemId:   i.ItemId(),
									Quantity: i.Quantity(),
									Slot:     i.Slot(),
								},
							})
						}

					}
				}
			}
		}
		if strings.Contains(include, "equipmentStatistics") {
			e, err := equipment.Processor(l, db).GetEquipmentForCharacter(uint32(characterId))
			if err != nil {
			} else {
				for _, equip := range e {
					es, err := statistics.Processor(l, db).GetEquipmentStatistics(equip.EquipmentId())
					if err != nil {
					} else {
						result.Included = append(result.Included, attributes.InventoryEquipmentStatisticsData{
							Id:   strconv.Itoa(int(es.Id())),
							Type: "com.atlas.cos.rest.attribute.EquipmentStatisticsAttributes",
							Attributes: attributes.InventoryEquipmentStatisticsAttributes{
								ItemId:        es.ItemId(),
								Strength:      es.Strength(),
								Dexterity:     es.Dexterity(),
								Intelligence:  es.Intelligence(),
								Luck:          es.Luck(),
								HP:            es.HP(),
								MP:            es.MP(),
								WeaponAttack:  es.WeaponAttack(),
								MagicAttack:   es.MagicAttack(),
								WeaponDefense: es.WeaponDefense(),
								MagicDefense:  es.MagicDefense(),
								Accuracy:      es.Accuracy(),
								Avoidability:  es.Avoidability(),
								Hands:         es.Hands(),
								Speed:         es.Speed(),
								Jump:          es.Jump(),
								Slots:         es.Slots(),
							},
						})
					}
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] %s.", err.Error())
		}
	}
}

func getInventoryItemType(inventoryType string) string {
	if inventoryType == inventory.ItemTypeEquip {
		return "com.atlas.cos.rest.attribute.EquipmentAttributes"
	} else {
		return "com.atlas.cos.rest.attribute.ItemAttributes"
	}
}

func GetInventoryForCharacter(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func CreateCharacterFromSeed(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		li := &attributes.CharacterSeedInputDataContainer{}
		err := attributes.FromJSON(li, r.Body)
		if err != nil {
			l.Println("[ERROR] deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		attr := li.Data.Attributes
		c, err := Processor(l, db).CreateFromSeed(attr.AccountId, attr.WorldId, attr.Name, attr.JobIndex, attr.Face, attr.Hair, attr.HairColor, attr.Skin, attr.Gender, attr.Top, attr.Bottom, attr.Shoes, attr.Weapon)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = &attributes.CharacterDataContainer{}
		result.Data = makeCharacterData(c)

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}

func GetSavedLocations(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theType := mux.Vars(r)["type"]

		var locations []*location.Model
		if theType != "" {
			locations, err = location.Processor(l, db).GetSavedLocationsByType(uint32(characterId), theType)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			locations, err = location.Processor(l, db).GetSavedLocations(uint32(characterId))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		var result attributes.LocationDataListContainer
		result.Data = make([]attributes.LocationData, 0)

		for _, loc := range locations {
			result.Data = append(result.Data, attributes.LocationData{
				Id:   strconv.Itoa(int(loc.Id())),
				Type: "com.atlas.cos.rest.attribute.LocationAttributes",
				Attributes: attributes.LocationAttributes{
					Type:     loc.Type(),
					MapId:    loc.MapId(),
					PortalId: loc.PortalId(),
				},
			})
		}

	}
}

func AddSavedLocation(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		li := &attributes.LocationInputDataContainer{}
		err = attributes.FromJSON(li, r.Body)
		if err != nil {
			l.Println("[ERROR] deserializing input", err)
			w.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		att := li.Data.Attributes
		err = location.Processor(l, db).AddSavedLocation(uint32(characterId), att.Type, att.MapId, att.PortalId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func GetCharacterDamage(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
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

func GetCharacterSkills(l *log.Logger, db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.Printf("[ERROR] unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sl, err := skill.Processor(l, db).GetSkills(uint32(characterId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var result = attributes.CharacterSkillDataListContainer{}
		result.Data = make([]attributes.CharacterSkillData, 0)
		for _, s := range sl {
			result.Data = append(result.Data, attributes.CharacterSkillData{
				Id:   strconv.Itoa(int(s.Id())),
				Type: "com.atlas.cos.rest.attribute.SkillAttributes",
				Attributes: attributes.CharacterSkillAttributes{
					Level:       s.Level(),
					MasterLevel: s.MasterLevel(),
					Expiration:  s.Expiration(),
				},
			})
		}

		w.WriteHeader(http.StatusOK)
		attributes.ToJSON(result, w)
	}
}
