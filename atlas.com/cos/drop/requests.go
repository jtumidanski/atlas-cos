package drop

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "drops/%d"
)

func requestById(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32) (*attributes.DropDataContainer, error) {
	return func(dropId uint32) (*attributes.DropDataContainer, error) {
		ar := &attributes.DropDataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(dropResource, dropId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
