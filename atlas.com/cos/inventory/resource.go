package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
	"atlas-cos/json"
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"atlas-cos/rest/resource"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func CreateItem(fl log.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fl.WithFields(log.Fields{"originator": "GetItemForCharacterByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inventoryType := mux.Vars(r)["type"]
		if inventoryType == "" {
			l.Errorf("Unable to retrieve requested inventory type.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.ToUpper(inventoryType) == TypeEquip {
			li := &attributes.EquipmentDataContainer{}
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

			itemId := li.Data.Attributes.ItemId
			ro, err := requests.EquipmentRegistry().Create(itemId)
			if err != nil {
				l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.", itemId, characterId)
				return
			}
			eid, err := strconv.Atoi(ro.Data.Id)
			if err != nil {
				l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.", itemId, characterId)
				return
			}

			err = equipment.GainItem(l, db)(uint32(characterId), itemId, uint32(eid))
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", characterId, itemId)
			}
		} else {
			li := &attributes.ItemDataContainer{}
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

			itemId := li.Data.Attributes.ItemId
			quantity := li.Data.Attributes.Quantity
			it, ok := GetByteFromName(inventoryType)
			if !ok {
				l.WithError(err).Errorf("Invalid inventory type supplied %s.", inventoryType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if quantity > 0 {
				err := item.GainItem(l, db)(uint32(characterId), it, itemId, uint32(quantity))
				if err != nil {
					l.WithError(err).Errorf("Unable to give character %d item %d.", characterId, itemId)
				}
			} else {
				err := item.LoseItem(l, db)(uint32(characterId), it, itemId, quantity)
				if err != nil {
					l.WithError(err).Errorf("Unable to take item %d from character %d.", itemId, characterId)
				}
			}

		}
	}
}

func GetItemForCharacterByType(l log.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetItemForCharacterByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		include := mux.Vars(r)["include"]

		inventoryType := mux.Vars(r)["type"]
		if inventoryType == "" {
			fl.Errorf("Unable to retrieve requested inventory type.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		slot, err := strconv.Atoi(mux.Vars(r)["slot"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse slot from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inv, err := GetInventory(fl, db)(uint32(characterId), inventoryType, FilterSlot(int16(slot)))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		prepareResult(fl, db, w)(uint32(characterId), inv, include)
	}
}

func GetItemsForCharacterByType(l log.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetItemForCharacterByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		include := mux.Vars(r)["include"]

		inventoryType := mux.Vars(r)["type"]
		if inventoryType == "" {
			fl.Errorf("Unable to retrieve requested inventory type.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		itemId, err := strconv.Atoi(mux.Vars(r)["itemId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse slot from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inv, err := GetInventory(fl, db)(uint32(characterId), inventoryType, FilterItemId(l, db)(uint32(itemId)))
		if err != nil {
			fl.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		prepareResult(fl, db, w)(uint32(characterId), inv, include)
	}
}

func GetInventoryForCharacterByType(l log.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetInventoryForCharacterByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		include := mux.Vars(r)["include"]

		inventoryType := mux.Vars(r)["type"]
		if inventoryType == "" {
			fl.Errorf("Unable to retrieve requested inventory type.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		inv, err := GetInventory(fl, db)(uint32(characterId), inventoryType)
		if err != nil {
			fl.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		prepareResult(fl, db, w)(uint32(characterId), inv, include)
	}
}

func prepareResult(l log.FieldLogger, db *gorm.DB, w http.ResponseWriter) func(characterId uint32, inv *Model, include string) {
	return func(characterId uint32, inv *Model, include string) {
		result := createInventoryDataContainer(inv)
		result.Data.Relationships.InventoryItems = createInventoryItemRelationships(inv)

		if strings.Contains(include, "inventoryItems") {
			result.Included = append(result.Included, createIncludedInventoryItems(l, db, characterId, inv)...)
		}
		if strings.Contains(include, "equipmentStatistics") {
			result.Included = append(result.Included, createIncludedEquipmentStatistics(l, db)(characterId, inv)...)
		}

		w.WriteHeader(http.StatusOK)
		err := json.ToJSON(result, w)
		if err != nil {
			l.WithError(err).Errorf("Writing response.")
		}
	}
}

func createIncludedEquipmentStatistics(l log.FieldLogger, db *gorm.DB) func(characterId uint32, inv *Model) []interface{} {
	return func(characterId uint32, inv *Model) []interface{} {
		var results = make([]interface{}, 0)
		e, err := equipment.GetEquipmentForCharacter(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", characterId)
			return results
		}

		for _, equip := range e {
			es, err := statistics.GetEquipmentStatistics(l)(equip.EquipmentId())
			if err != nil {
				l.WithError(err).Errorf("Retrieving equipment %d statistics for character %d.", equip.EquipmentId(), characterId)
			} else {
				ok := false
				for _, e := range inv.Items() {
					if e.Id() == equip.Id() {
						ok = true
						break
					}
				}
				if ok {
					results = append(results, createEquipmentStatisticsData(es))
				}
			}
		}
		return results
	}
}

func createIncludedInventoryItems(fl log.FieldLogger, db *gorm.DB, characterId uint32, inv *Model) []interface{} {
	var results = make([]interface{}, 0)
	for _, inventoryItem := range inv.Items() {
		if inventoryItem.Type() == ItemTypeEquip {
			e, err := equipment.GetEquipmentById(fl, db)(inventoryItem.Id())
			if err != nil {
				fl.WithError(err).Errorf("Unable to retrieve equipment %d for character %d.", inventoryItem.Id(), characterId)
			} else {
				results = append(results, createEquipmentData(e))
			}
		} else {
			i, err := item.GetItemById(fl, db)(inventoryItem.Id())
			if err != nil {
				fl.WithError(err).Errorf("Unable to retrieve item %d for character %d.", inventoryItem.Id(), characterId)
			} else {
				results = append(results, createItemData(i))
			}

		}
	}
	return results
}

func createInventoryItemRelationships(inv *Model) []attributes.Relationship {
	var results = make([]attributes.Relationship, 0)
	for _, i := range inv.Items() {
		var inventoryItem = attributes.Relationship{
			Id:   strconv.Itoa(int(i.Id())),
			Type: getInventoryItemType(i.Type()),
		}
		results = append(results, inventoryItem)
	}
	return results
}

func createInventoryDataContainer(inv *Model) *attributes.InventoryDataContainer {
	return &attributes.InventoryDataContainer{Data: createInventoryData(inv)}
}

func createInventoryData(inv *Model) attributes.InventoryData {
	return attributes.InventoryData{
		Id:   strconv.Itoa(int(inv.Id())),
		Type: "com.atlas.cos.rest.attribute.InventoryAttributes",
		Attributes: attributes.InventoryAttributes{
			Type:     inv.Type(),
			Capacity: inv.Capacity(),
		},
	}
}

func createEquipmentData(e *equipment.Model) attributes.InventoryEquipmentData {
	return attributes.InventoryEquipmentData{
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
	}
}

func createItemData(i *item.Model) attributes.InventoryItemData {
	return attributes.InventoryItemData{
		Id:   strconv.Itoa(int(i.Id())),
		Type: "com.atlas.cos.rest.attribute.ItemAttributes",
		Attributes: attributes.InventoryItemAttributes{
			ItemId:   i.ItemId(),
			Quantity: int32(i.Quantity()),
			Slot:     i.Slot(),
		},
	}
}

func createEquipmentStatisticsData(es *statistics.Model) attributes.InventoryEquipmentStatisticsData {
	return attributes.InventoryEquipmentStatisticsData{
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
	}
}

func getInventoryItemType(inventoryType string) string {
	if inventoryType == TypeEquip {
		return "com.atlas.cos.rest.attribute.EquipmentAttributes"
	} else {
		return "com.atlas.cos.rest.attribute.ItemAttributes"
	}
}

func GetItemsForCharacter(l log.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fl := l.WithFields(log.Fields{"originator": "GetItemForCharacterByType", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		itemId, err := strconv.Atoi(mux.Vars(r)["itemId"])
		if err != nil {
			fl.WithError(err).Errorf("Unable to properly parse itemId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result := &ItemListDataContainer{}

		types := []string{TypeEquip, TypeUse, TypeSetup, TypeETC, TypeCash}
		for _, t := range types {
			inv, err := GetInventory(fl, db)(uint32(characterId), t, FilterItemId(l, db)(uint32(itemId)))
			if err != nil {
				fl.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, t)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, i := range inv.Items() {
				quantity := uint32(1)
				if i.ItemType() == ItemTypeItem {
					ii, err := item.GetItemById(l, db)(i.Id())
					if err != nil {
						l.WithError(err).Errorf("Unable to lookup item by id %d.", i.Id())
						continue
					}
					quantity = ii.Quantity()
				}

				result.Data = append(result.Data, ItemDataBody{
					Id:   strconv.Itoa(int(i.Id())),
					Type: i.Type(),
					Attributes: ItemAttributes{
						InventoryType: inv.Type(),
						Slot:          i.Slot(),
						Quantity:      quantity,
					},
				})
			}
		}

		w.WriteHeader(http.StatusOK)
		err = json.ToJSON(result, w)
		if err != nil {
			l.WithError(err).Errorf("Writing response.")
		}
	}
}
