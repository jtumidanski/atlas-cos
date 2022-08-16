package inventory

import "errors"

func makeInventory(e entity) (Model, error) {
	t, ok := GetTypeFromByte(e.InventoryType)
	if !ok {
		return Model{}, errors.New("invalid type")
	}

	return Model{
		id:            e.ID,
		inventoryType: t,
		capacity:      e.Capacity,
	}, nil
}
