package equipment

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	equipmentServicePrefix string = "/ms/eso/"
	equipmentService              = requests.BaseRequest + equipmentServicePrefix
	equipmentResource             = equipmentService + "equipment"
	equipResource                 = equipmentResource + "/%d"
)

func requestCreate(itemId uint32) (*attributes.EquipmentDataContainer, error) {
	input := &attributes.EquipmentDataContainer{
		Data: attributes.EquipmentData{
			Id:   "0",
			Type: "com.atlas.eso.attribute.EquipmentAttributes",
			Attributes: attributes.EquipmentAttributes{
				ItemId: itemId,
			},
		}}
	resp, err := requests.Post(fmt.Sprintf(equipmentResource), input)
	if err != nil {
		return nil, err
	}

	ro := &attributes.EquipmentDataContainer{}
	err = requests.ProcessResponse(resp, ro)
	if err != nil {
		return nil, err
	}
	return ro, nil
}

func requestById(l logrus.FieldLogger) func(equipmentId uint32) (*attributes.EquipmentDataContainer, error) {
	return func(equipmentId uint32) (*attributes.EquipmentDataContainer, error) {
		ar := &attributes.EquipmentDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(equipResource, equipmentId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
