package item

import "gorm.io/gorm"

type EntityUpdateFunction func(e *entity)

func CreateItemForCharacter(db *gorm.DB, characterId uint32, inventoryType byte, itemId uint32, quantity uint32, slot int16) (*Model, error) {
	e := &entity{
		CharacterId:   characterId,
		InventoryType: inventoryType,
		ItemId:        itemId,
		Quantity:      quantity,
		Slot:          slot,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeItem(e), nil
}

func Update(db *gorm.DB, uniqueId uint32, modifiers ...EntityUpdateFunction) error {
	e := &entity{}
	for _, modifier := range modifiers {
		modifier(e)
	}
	return db.Model(&entity{Id: uniqueId}).Updates(e).Error
}

func SetQuantity(quantity uint32) EntityUpdateFunction {
	return func(e *entity) {
		e.Quantity = quantity
	}
}
