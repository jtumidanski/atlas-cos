package item

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	itemInformationServicePrefix     string = "/ms/iis/"
	itemInformationService                  = requests.BaseRequest + itemInformationServicePrefix
	itemInformationEquipmentResource        = itemInformationService + "equipment"
	itemInformationEquipResource            = itemInformationEquipmentResource + "/%d"
	itemInformationEquipSlotResource        = itemInformationEquipResource + "/slots"
)

func requestEquipmentSlotDestination(l logrus.FieldLogger) func(itemId uint32) (*attributes.EquipmentSlotDataListContainer, error) {
	return func(itemId uint32) (*attributes.EquipmentSlotDataListContainer, error) {
		ar := &attributes.EquipmentSlotDataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(itemInformationEquipSlotResource, itemId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
