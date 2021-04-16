package item

import "gorm.io/gorm"

type EntityUpdateFunction func(e entity)

func CreateItemForCharacter(db *gorm.DB, characterId uint32, inventoryType byte, itemId uint32, quantity uint32, slot int16) (*Model, error) {
	e := &entity{
		characterId:   characterId,
		inventoryType: inventoryType,
		itemId:        itemId,
		quantity:      quantity,
		slot:          slot,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeItem(e), nil
}

func Update(db *gorm.DB, uniqueId uint32, modifiers ...EntityUpdateFunction) error {
	c := entity{id: uniqueId}
	err := db.Where(&c).First(&c).Error
	if err != nil {
		return err
	}

	for _, modifier := range modifiers {
		modifier(c)
	}

	err = db.Save(&c).Error
	return err
}

func SetQuantity(quantity uint32) EntityUpdateFunction {
	return func(e entity) {
		e.quantity = quantity
	}
}