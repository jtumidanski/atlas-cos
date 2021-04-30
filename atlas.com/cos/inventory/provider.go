package inventory

import "gorm.io/gorm"

func Get(db *gorm.DB, characterId uint32, inventoryType int8) (*Model, error) {
	var result entity
	err := db.Where(&entity{CharacterId: characterId, InventoryType: inventoryType}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return makeInventory(&result), nil
}