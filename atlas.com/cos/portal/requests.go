package portal

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	portalsResource                    = mapsResource + "%d/portals"
	portalResource                     = portalsResource + "/%d"
)

func requestById(mapId uint32, portalId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(portalResource, mapId, portalId))
}

func requestInMap(mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(portalsResource, mapId))
}
