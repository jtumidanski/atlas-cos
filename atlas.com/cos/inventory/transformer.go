package inventory

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