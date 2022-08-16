package item

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	itemInformationServicePrefix     string = "/ms/iis/"
	itemInformationService                  = requests.BaseRequest + itemInformationServicePrefix
	itemInformationEquipmentResource        = itemInformationService + "equipment"
	itemInformationEquipResource            = itemInformationEquipmentResource + "/%d"
	itemInformationEquipSlotResource        = itemInformationEquipResource + "/slots"
)

func requestEquipmentSlotDestination(itemId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(itemInformationEquipSlotResource, itemId))
}
