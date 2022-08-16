package item

import (
	"errors"
	"gorm.io/gorm"
)

const TypeEquipment = "equipment"
const TypeItem = "item"

func createItem(db *gorm.DB, inventoryId uint32, itemId uint32, quantity uint32, slot int16) (ItemModel, error) {
	var im ItemModel
	txError := db.Transaction(func(tx *gorm.DB) error {
		ei := &entityItem{
			Quantity: quantity,
		}
		err := db.Create(ei).Error
		if err != nil {
			return err
		}

		eii := &entityInventoryItem{
			InventoryId: inventoryId,
			ItemId:      itemId,
			Slot:        slot,
			Type:        TypeItem,
			ReferenceId: ei.ID,
		}
		err = db.Create(eii).Error
		if err != nil {
			return err
		}
		im, err = makeItem(*eii, *ei)
		if err != nil {
			return err
		}
		return nil
	})
	return im, txError
}

func createEquipment(db *gorm.DB, inventoryId uint32, itemId uint32, slot int16, equipmentId uint32) (EquipmentModel, error) {
	e := &entityInventoryItem{
		InventoryId: inventoryId,
		ItemId:      itemId,
		Slot:        slot,
		Type:        TypeEquipment,
		ReferenceId: equipmentId,
	}

	err := db.Create(e).Error
	if err != nil {
		return EquipmentModel{}, err
	}
	return makeEquipment(*e)
}

func makeModel(db *gorm.DB) func(item entityInventoryItem) (Model, error) {
	return func(ii entityInventoryItem) (Model, error) {
		if ii.Type == TypeEquipment {
			return makeEquipment(ii)
		}
		if ii.Type == TypeItem {
			ei, err := getItemAttributes(ii.ReferenceId)(db)()
			if err != nil {
				return ItemModel{}, err
			}

			return makeItem(ii, ei)
		}
		return nil, errors.New("unsupported item type")
	}
}

func makeItem(ii entityInventoryItem, i entityItem) (ItemModel, error) {
	return ItemModel{
		id:          ii.ID,
		inventoryId: ii.InventoryId,
		itemId:      ii.ItemId,
		slot:        ii.Slot,
		quantity:    i.Quantity,
		referenceId: ii.ReferenceId,
	}, nil
}

func makeEquipment(item entityInventoryItem) (EquipmentModel, error) {
	return EquipmentModel{
		id:          item.ID,
		inventoryId: item.InventoryId,
		itemId:      item.ItemId,
		slot:        item.Slot,
		equipmentId: item.ReferenceId,
	}, nil
}

func remove(db *gorm.DB, inventoryId uint32, id uint32) error {
	return db.Delete(&entityInventoryItem{InventoryId: inventoryId, ID: id}).Error
}

func updateQuantity(db *gorm.DB, id uint32, amount uint32) error {
	return db.Model(&entityItem{ID: id}).Select("Quantity").Updates(&entityItem{Quantity: amount}).Error
}

func updateSlot(db *gorm.DB, id uint32, slot int16) error {
	return db.Model(&entityInventoryItem{ID: id}).Select("Slot").Updates(&entityInventoryItem{Slot: slot}).Error
}
