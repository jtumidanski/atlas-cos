package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	equipmentServicePrefix string = "/ms/eso/"
	equipmentService              = baseRequest + equipmentServicePrefix
	equipmentResource             = equipmentService + "equipment"
	equipResource                 = equipmentResource + "/%d"
)

var EquipmentRegistry = func() *equipmentRegistry {
	return &equipmentRegistry{}
}

type equipmentRegistry struct {
}

func (e *equipmentRegistry) Create(itemId uint32) (*attributes.EquipmentDataContainer, error) {
	input := &attributes.EquipmentDataContainer{
		Data: attributes.EquipmentData{
			Id:   "0",
			Type: "com.atlas.eso.attribute.EquipmentAttributes",
			Attributes: attributes.EquipmentAttributes{
				ItemId: itemId,
			},
		}}
	resp, err := post(fmt.Sprintf(equipmentResource), input)
	if err != nil {
		return nil, err
	}

	ro := &attributes.EquipmentDataContainer{}
	err = processResponse(resp, ro)
	if err != nil {
		return nil, err
	}
	return ro, nil
}

func (e *equipmentRegistry) GetById(equipmentId uint32) (*attributes.EquipmentDataContainer, error) {
	ar := &attributes.EquipmentDataContainer{}
	err := get(fmt.Sprintf(equipResource, equipmentId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
