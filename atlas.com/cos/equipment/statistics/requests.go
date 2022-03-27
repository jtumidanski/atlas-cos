package statistics

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	equipmentServicePrefix string = "/ms/eso/"
	equipmentService              = requests.BaseRequest + equipmentServicePrefix
	equipmentResource             = equipmentService + "equipment"
	equipResource                 = equipmentResource + "/%d"
)

func requestCreate(itemId uint32) requests.PostRequest[attributes] {
	input := &EquipmentDataContainer{
		Data: EquipmentData{
			Id:   "0",
			Type: "com.atlas.eso.attribute.EquipmentAttributes",
			Attributes: attributes{
				ItemId: itemId,
			},
		},
	}
	return requests.MakePostRequest[attributes](fmt.Sprintf(equipmentResource), input)
}

func requestById(equipmentId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(equipResource, equipmentId))
}
