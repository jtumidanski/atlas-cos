package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
	"atlas-cos/rest/attributes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetInventoryForCharacterByType(l *log.Logger, db *gorm.DB) http.HandlerFunc {
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

		inv, err := Processor(l, db).GetInventoryByType(uint32(characterId), inventoryType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := createInventoryDataContainer(inv)
		result.Data.Relationships.InventoryItems = createInventoryItemRelationships(inv)

		if strings.Contains(include, "inventoryItems") {
			result.Included = append(result.Included, createIncludedInventoryItems(l, db, uint32(characterId), inventoryType))
		}
		if strings.Contains(include, "equipmentStatistics") {
			result.Included = append(result.Included, createIncludedEquipmentStatistics(l, db, uint32(characterId)))
		}

		w.WriteHeader(http.StatusOK)
		err = attributes.ToJSON(result, w)
		if err != nil {
			l.Printf("[ERROR] %s.", err.Error())
		}
	}
}

func GetInventoryForCharacter(l *log.Logger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO
	}
}

func createIncludedEquipmentStatistics(l *log.Logger, db *gorm.DB, characterId uint32) []interface{} {
	var results = make([]interface{}, 0)
	e, err := equipment.Processor(l, db).GetEquipmentForCharacter(characterId)
	if err != nil {
		l.Printf("[ERROR] unable to retrieve equipment for character %d.", characterId)
		return results
	}

	for _, equip := range e {
		es, err := statistics.Processor(l, db).GetEquipmentStatistics(equip.EquipmentId())
		if err != nil {
			l.Printf("[ERROR] retrieving equipment %d statistics for character %d.", equip.EquipmentId(), characterId)
		} else {
			results = append(results, createEquipmentStatisticsData(es))
		}
	}
	return results
}

func createIncludedInventoryItems(l *log.Logger, db *gorm.DB, characterId uint32, inventoryType string) []interface{} {
	var results = make([]interface{}, 0)
	inv, err := Processor(l, db).GetInventoryByType(characterId, inventoryType)
	if err != nil {
		l.Printf("[ERROR] unable to retrieve inventory items for character %d.", characterId)
		return results
	}
	for _, inventoryItem := range inv.Items() {
		if inventoryItem.Type() == ItemTypeEquip {
			e, err := equipment.Processor(l, db).GetEquipmentById(inventoryItem.Id())
			if err != nil {
				l.Printf("[ERROR] unable to retrieve equipment %d for character %d.", inventoryItem.Id(), characterId)
			} else {
				results = append(results, createEquipmentData(e))
			}
		} else {
			i, err := item.Processor(l, db).GetItemById(inventoryItem.Id())
			if err != nil {
				l.Printf("[ERROR] unable to retrieve item %d for character %d.", inventoryItem.Id(), characterId)
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
			Quantity: i.Quantity(),
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
	if inventoryType == ItemTypeEquip {
		return "com.atlas.cos.rest.attribute.EquipmentAttributes"
	} else {
		return "com.atlas.cos.rest.attribute.ItemAttributes"
	}
}
