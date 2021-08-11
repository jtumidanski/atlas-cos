package portal

import (
	"atlas-cos/rest/attributes"
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	portalsResource                    = mapsResource + "%d/portals"
	portalsByName                      = portalsResource + "?name=%s"
	portalResource                     = portalsResource + "/%d"
)

func requestPortalById(l logrus.FieldLogger) func(mapId uint32, portalId uint32) (*attributes.PortalDataContainer, error) {
	return func(mapId uint32, portalId uint32) (*attributes.PortalDataContainer, error) {
		ar := &attributes.PortalDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(portalResource, mapId, portalId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestPortals(l logrus.FieldLogger) func(mapId uint32) (*attributes.PortalDataContainer, error) {
	return func(mapId uint32) (*attributes.PortalDataContainer, error) {
		ar := &attributes.PortalDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(portalsResource, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
