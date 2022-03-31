package item

import "gorm.io/gorm"

type EntityUpdateFunction func() ([]string, func(e *entity))

func createItemForCharacter(db *gorm.DB, characterId uint32, inventoryType int8, itemId uint32, quantity uint32, slot int16) (Model, error) {
	e := &entity{
		CharacterId:   characterId,
		InventoryType: inventoryType,
		ItemId:        itemId,
		Quantity:      quantity,
		Slot:          slot,
	}

	err := db.Create(e).Error
	if err != nil {
		return Model{}, err
	}
	return makeItem(*e)
}

func remove(db *gorm.DB, uniqueId uint32) error {
	return db.Delete(&entity{Id: uniqueId}).Error
}

func update(db *gorm.DB, uniqueId uint32, modifiers ...EntityUpdateFunction) error {
	e := &entity{}
	var columns []string
	for _, modifier := range modifiers {
		c, u := modifier()
		columns = append(columns, c...)
		u(e)
	}
	return db.Model(&entity{Id: uniqueId}).Select(columns).Updates(e).Error
}

func SetQuantity(quantity uint32) EntityUpdateFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Quantity"}, func(e *entity) {
			e.Quantity = quantity
		}
	}
}

func SetSlot(slot int16) EntityUpdateFunction {
	return func() ([]string, func(e *entity)) {
		return []string{"Slot"}, func(e *entity) {
			e.Slot = slot
		}
	}
}
