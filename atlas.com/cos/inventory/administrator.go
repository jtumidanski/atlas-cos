package inventory

import "gorm.io/gorm"

func Create(db *gorm.DB, characterId uint32, inventoryType byte, capacity uint32) (*Model, error) {
	e := &entity{
		CharacterId:   characterId,
		InventoryType: inventoryType,
		Capacity:      capacity,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeInventory(e), nil
}

func makeInventory(e *entity) *Model {
	t, ok := GetTypeFromByte(e.InventoryType)
	if !ok {
		return nil
	}

	return &Model{
		id:            e.InventoryType,
		inventoryType: t,
		capacity:      e.Capacity,
	}
}
