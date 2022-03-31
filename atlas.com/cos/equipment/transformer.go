package equipment

func makeEquipment(e entity) (Model, error) {
	return Model{
		id:          e.Id,
		characterId: e.CharacterId,
		equipmentId: e.EquipmentId,
		slot:        e.Slot,
	}, nil
}
