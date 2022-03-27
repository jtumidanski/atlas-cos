package inventory

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/item"
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
	"strings"
)

const (
	GetItemsForCharacter           = "get_items_for_character"
	GetItemForCharacterByType      = "get_item_for_character_by_type"
	GetItemsForCharacterByType     = "get_items_for_character_by_type"
	GetInventoryForCharacterByType = "get_inventory_for_character_by_type"
	CreateItem                     = "create_item"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	r := router.PathPrefix("/characters").Subrouter()
	r.HandleFunc("/{characterId}/items", registerGetItemsForCharacter(l, db)).Methods(http.MethodGet).Queries("itemId", "{itemId}")
	r.HandleFunc("/{characterId}/inventories", registerGetItemForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}", "slot", "{slot}")
	r.HandleFunc("/{characterId}/inventories", registerGetItemsForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}", "itemId", "{itemId}")
	r.HandleFunc("/{characterId}/inventories", registerGetInventoryForCharacterByType(l, db)).Methods(http.MethodGet).Queries("include", "{include}", "type", "{type}")
	r.HandleFunc("/{characterId}/inventories/{type}/items", registerCreateItem(l, db)).Methods(http.MethodPost)
}

func registerCreateItem(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(CreateItem, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": CreateItem, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleCreateItem(fl, db)(span)(characterId)
		})
	})
}

func registerGetInventoryForCharacterByType(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetInventoryForCharacterByType, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": GetInventoryForCharacterByType, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetInventoryForCharacterByType(fl, db)(span)(characterId)
		})
	})
}

func registerGetItemsForCharacterByType(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetItemsForCharacterByType, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": GetItemsForCharacterByType, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetItemsForCharacterByType(fl, db)(span)(characterId)
		})
	})
}

func registerGetItemForCharacterByType(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetItemForCharacterByType, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": GetItemForCharacterByType, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetItemForCharacterByType(fl, db)(span)(characterId)
		})
	})
}

func registerGetItemsForCharacter(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetItemsForCharacter, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": GetItemsForCharacter, "type": "rest_handler"})
		return parseCharacterId(fl, func(characterId uint32) http.HandlerFunc {
			return handleGetItemsForCharacter(l, db)(span)(characterId)
		})
	})
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

func handleCreateItem(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				inventoryType := mux.Vars(r)["type"]
				if inventoryType == "" {
					l.Errorf("Unable to retrieve requested inventory type.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if strings.ToUpper(inventoryType) == TypeEquip {
					li := &statistics.EquipmentDataContainer{}
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
					eid, err := statistics.Create(l, span)(itemId)
					if err != nil {
						l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.", itemId, characterId)
						return
					}

					err = GainEquipment(l, db, span)(characterId, itemId, eid)
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
						err = GainItem(l, db, span)(characterId, it, itemId, uint32(quantity))
						if err != nil {
							l.WithError(err).Errorf("Unable to give character %d item %d.", characterId, itemId)
						}
					} else {
						err = LoseItem(l, db, span)(characterId, it, itemId, quantity)
						if err != nil {
							l.WithError(err).Errorf("Unable to take item %d from character %d.", itemId, characterId)
						}
					}

				}
			}
		}
	}
}

func handleGetItemForCharacterByType(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				include := mux.Vars(r)["include"]

				inventoryType := mux.Vars(r)["type"]
				if inventoryType == "" {
					l.Errorf("Unable to retrieve requested inventory type.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				slot, err := strconv.Atoi(mux.Vars(r)["slot"])
				if err != nil {
					l.WithError(err).Errorf("Unable to properly parse slot from path.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				inv, err := GetInventory(l, db)(characterId, inventoryType, FilterSlot(int16(slot)))
				if err != nil {
					l.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				prepareResult(l, db, span, w)(characterId, inv, include)
			}
		}
	}
}

func handleGetItemsForCharacterByType(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				include := mux.Vars(r)["include"]
				inventoryType := mux.Vars(r)["type"]

				if inventoryType == "" {
					l.Errorf("Unable to retrieve requested inventory type.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				itemId, err := strconv.Atoi(mux.Vars(r)["itemId"])
				if err != nil {
					l.WithError(err).Errorf("Unable to properly parse slot from path.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				inv, err := GetInventory(l, db)(characterId, inventoryType, FilterItemId(l, db, span)(uint32(itemId)))
				if err != nil {
					l.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				prepareResult(l, db, span, w)(characterId, inv, include)
			}
		}
	}
}

func handleGetInventoryForCharacterByType(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				include := mux.Vars(r)["include"]
				inventoryType := mux.Vars(r)["type"]
				if inventoryType == "" {
					l.Errorf("Unable to retrieve requested inventory type.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				inv, err := GetInventory(l, db)(characterId, inventoryType)
				if err != nil {
					l.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, inventoryType)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				prepareResult(l, db, span, w)(characterId, inv, include)
			}
		}
	}
}

func prepareResult(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span, w http.ResponseWriter) func(characterId uint32, inv *Model, include string) {
	return func(characterId uint32, inv *Model, include string) {
		result := createInventoryDataContainer(inv)
		result.Data.Relationships.InventoryItems = createInventoryItemRelationships(inv)

		if strings.Contains(include, "inventoryItems") {
			result.Included = append(result.Included, createIncludedInventoryItems(l, db, characterId, inv)...)
		}
		if strings.Contains(include, "equipmentStatistics") {
			result.Included = append(result.Included, createIncludedEquipmentStatistics(l, db, span)(characterId, inv)...)
		}

		w.WriteHeader(http.StatusOK)
		err := json.ToJSON(result, w)
		if err != nil {
			l.WithError(err).Errorf("Writing response.")
		}
	}
}

func createIncludedEquipmentStatistics(l logrus.FieldLogger, db *gorm.DB, span opentracing.Span) func(characterId uint32, inv *Model) []interface{} {
	return func(characterId uint32, inv *Model) []interface{} {
		var results = make([]interface{}, 0)
		e, err := equipment.GetEquipmentForCharacter(l, db)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve equipment for character %d.", characterId)
			return results
		}

		for _, equip := range e {
			es, err := statistics.GetEquipmentStatistics(l, span)(equip.EquipmentId())
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

func createIncludedInventoryItems(fl logrus.FieldLogger, db *gorm.DB, characterId uint32, inv *Model) []interface{} {
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

func handleGetItemsForCharacter(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) http.HandlerFunc {
		return func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				itemId, err := strconv.Atoi(mux.Vars(r)["itemId"])
				if err != nil {
					l.WithError(err).Errorf("Unable to properly parse itemId from path.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				result := &ItemListDataContainer{}

				types := []string{TypeEquip, TypeUse, TypeSetup, TypeETC, TypeCash}
				for _, t := range types {
					inv, err := GetInventory(l, db)(characterId, t, FilterItemId(l, db, span)(uint32(itemId)))
					if err != nil {
						l.WithError(err).Errorf("Unable to get inventory for character %d by type %s.", characterId, t)
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
	}
}
