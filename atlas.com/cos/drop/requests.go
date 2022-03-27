package drop

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "drops/%d"
)

func requestById(dropId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(dropResource, dropId))
}
