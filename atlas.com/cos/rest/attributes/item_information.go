package attributes

type EquipmentSlotDataContainer struct {
	Data EquipmentSlotData `json:"data"`
}

type EquipmentSlotDataListContainer struct {
	Data []EquipmentSlotData `json:"data"`
}

type EquipmentSlotData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes EquipmentSlotAttributes `json:"attributes"`
}

type EquipmentSlotAttributes struct {
	Name string `json:"name"`
	WZ   string `json:"wz"`
	Slot int16  `json:"slot"`
}
