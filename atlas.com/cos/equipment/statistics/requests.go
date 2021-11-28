package statistics

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	equipmentServicePrefix string = "/ms/eso/"
	equipmentService              = requests.BaseRequest + equipmentServicePrefix
	equipmentResource             = equipmentService + "equipment"
	equipResource                 = equipmentResource + "/%d"
)

func requestCreate(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) (*attributes.EquipmentDataContainer, error) {
	return func(itemId uint32) (*attributes.EquipmentDataContainer, error) {
		input := &attributes.EquipmentDataContainer{
			Data: attributes.EquipmentData{
				Id:   "0",
				Type: "com.atlas.eso.attribute.EquipmentAttributes",
				Attributes: attributes.EquipmentAttributes{
					ItemId: itemId,
				},
			}}

		ro := &attributes.EquipmentDataContainer{}
		err := requests.Post(l, span)(fmt.Sprintf(equipmentResource), input, ro, &requests.ErrorListDataContainer{})
		if err != nil {
			return nil, err
		}
		return ro, nil
	}
}

func requestById(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) (*attributes.EquipmentDataContainer, error) {
	return func(equipmentId uint32) (*attributes.EquipmentDataContainer, error) {
		ar := &attributes.EquipmentDataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(equipResource, equipmentId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
