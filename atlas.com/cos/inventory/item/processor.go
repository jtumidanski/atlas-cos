package item

import (
	"atlas-cos/database"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sort"
)

var characterCreationItems = []uint32{
	1302000, 1312004, 1322005, 1442079, // weapons
	1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167, // top
	1060002, 1060006, 1061002, 1061008, 1062115, // bottom
	1072001, 1072005, 1072037, 1072038, 1072383, // shoes
	30000, 30010, 30020, 30030, 31000, 31040, 31050, // hair
	20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100, //face
}

func invalidCharacterCreationItem(itemId uint32) bool {
	for _, v := range characterCreationItems {
		if itemId == v {
			return false
		}
	}
	return true
}

func CreateEquipment(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryId uint32, itemId uint32, equipmentId uint32, characterCreation bool) (EquipmentModel, error) {
	return func(characterId uint32, inventoryId uint32, itemId uint32, equipmentId uint32, characterCreation bool) (EquipmentModel, error) {
		if characterCreation {
			if invalidCharacterCreationItem(itemId) {
				l.Errorf("Received a request to create an item %d for character %d which is not valid for character creation. This is usually a hack.")
				return EquipmentModel{}, errors.New("not valid item for character creation")
			}
		}

		nextOpen, err := GetNextFreeSlot(l, db)(inventoryId)
		if err != nil {
			nextOpen = 1
		}

		eq, err := createEquipment(db, inventoryId, itemId, nextOpen, equipmentId)
		if err != nil {
			l.Errorf("Persisting equipment %d association for character %d in Slot %d.", equipmentId, characterId, nextOpen)
			return EquipmentModel{}, err
		}
		return eq, nil
	}
}

func CreateItem(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, inventoryId uint32, inventoryType int8, itemId uint32, quantity uint32) (Model, error) {
	return func(characterId uint32, inventoryId uint32, inventoryType int8, itemId uint32, quantity uint32) (Model, error) {
		slot, err := GetNextFreeSlot(l, db)(inventoryId)
		if err != nil {
			return ItemModel{}, err
		}
		return createItem(db, inventoryId, itemId, quantity, slot)
	}
}

func GetEquipment(_ logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32) ([]EquipmentModel, error) {
	return func(inventoryId uint32) ([]EquipmentModel, error) {
		return database.ModelSliceProvider[EquipmentModel, entityInventoryItem](db)(getByInventory(inventoryId), makeEquipment)()
	}
}

func GetByInventory(_ logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32) ([]Model, error) {
	return func(inventoryId uint32) ([]Model, error) {
		return database.ModelSliceProvider[Model, entityInventoryItem](db)(getByInventory(inventoryId), makeModel(db))()
	}
}

func GetBySlot(_ logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32, slot int16) (Model, error) {
	return func(inventoryId uint32, slot int16) (Model, error) {
		return database.ModelProvider[Model, entityInventoryItem](db)(getBySlot(inventoryId, slot), makeModel(db))()
	}
}

func GetEquippedItemBySlot(_ logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32, slot int16) (EquipmentModel, error) {
	return func(inventoryId uint32, slot int16) (EquipmentModel, error) {
		return database.ModelProvider[EquipmentModel, entityInventoryItem](db)(getBySlot(inventoryId, slot), makeEquipment)()
	}
}

func GetById(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) (Model, error) {
	return func(id uint32) (Model, error) {
		return database.ModelProvider[Model, entityInventoryItem](db)(getById(id), makeModel(db))()
	}
}

func GetByEquipmentId(_ logrus.FieldLogger, db *gorm.DB) func(id uint32) (EquipmentModel, error) {
	return func(id uint32) (EquipmentModel, error) {
		return database.ModelProvider[EquipmentModel, entityInventoryItem](db)(getByEquipmentId(id), makeEquipment)()
	}
}

func GetByItemId(_ logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32, itemId uint32) ([]Model, error) {
	return func(inventoryId uint32, itemId uint32) ([]Model, error) {
		return database.ModelSliceProvider[Model, entityInventoryItem](db)(getForCharacter(inventoryId, itemId), makeModel(db))()
	}
}

func UpdateSlot(_ logrus.FieldLogger, db *gorm.DB) func(id uint32, slot int16) error {
	return func(id uint32, slot int16) error {
		return updateSlot(db, id, slot)
	}
}

func UpdateQuantity(l logrus.FieldLogger, db *gorm.DB) func(id uint32, quantity uint32) error {
	return func(id uint32, quantity uint32) error {
		i, err := GetById(l, db)(id)
		if err != nil {
			return err
		}
		if val, ok := i.(*ItemModel); ok {
			return updateQuantity(db, val.ReferenceId(), quantity)
		}
		return nil
	}
}

func MaxInSlot() uint32 {
	return 200
}

func GetEquipmentSlotDestination(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) ([]int16, error) {
	return func(itemId uint32) ([]int16, error) {
		r, err := requestEquipmentSlotDestination(itemId)(l, span)
		if err != nil {
			return nil, err
		}

		var slots = make([]int16, 0)
		for _, data := range r.DataList() {
			attr := data.Attributes
			slots = append(slots, attr.Slot)
		}
		return slots, nil
	}
}

func GetNextFreeSlot(l logrus.FieldLogger, db *gorm.DB) func(inventoryId uint32) (int16, error) {
	return func(inventoryId uint32) (int16, error) {
		es, err := GetByInventory(l, db)(inventoryId)
		if err != nil {
			return 1, err
		}
		if len(es) == 0 {
			return 1, nil
		}

		sort.Slice(es, func(i, j int) bool {
			return es[i].Slot() < es[j].Slot()
		})
		return minFreeSlot(es), nil
	}
}

func RemoveItem(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32, id uint32) error {
	return func(characterId uint32, id uint32) error {
		return remove(db, characterId, id)
	}
}
