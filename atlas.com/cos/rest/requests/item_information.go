package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	itemInformationServicePrefix     string = "/ms/iis/"
	itemInformationService                  = BaseRequest + itemInformationServicePrefix
	itemInformationEquipmentResource        = itemInformationService + "equipment"
	itemInformationEquipResource            = itemInformationEquipmentResource + "/%d"
	itemInformationEquipSlotResource        = itemInformationEquipResource + "/slots"
)

var ItemInformationRegistry = func() *itemInformationRegistry {
	return &itemInformationRegistry{}
}

type itemInformationRegistry struct {
}

func (r itemInformationRegistry) GetEquipmentSlotDestination(itemId uint32) (*attributes.EquipmentSlotDataListContainer, error) {
	ar := &attributes.EquipmentSlotDataListContainer{}
	err := get(fmt.Sprintf(itemInformationEquipSlotResource, itemId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}