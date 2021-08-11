package drop

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "drops/%d"
)

func requestById(l logrus.FieldLogger) func(dropId uint32) (*attributes.DropDataContainer, error) {
	return func(dropId uint32) (*attributes.DropDataContainer, error) {
		ar := &attributes.DropDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(dropResource, dropId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
