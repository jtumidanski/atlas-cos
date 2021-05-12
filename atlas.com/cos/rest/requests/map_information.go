package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	portalsResource                    = mapsResource + "%d/portals"
	portalsByName                      = portalsResource + "?name=%s"
	portalResource                     = portalsResource + "/%d"
)

var MapInformation = func() *mapInformation {
	return &mapInformation{}
}

type mapInformation struct {
}

func (m *mapInformation) GetPortalById(mapId uint32, portalId uint32) (*attributes.PortalDataContainer, error) {
	ar := &attributes.PortalDataContainer{}
	err := get(fmt.Sprintf(portalResource, mapId, portalId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (m *mapInformation) GetPortalByName(mapId uint32, portalName string) (*attributes.PortalDataContainer, error) {
	ar := &attributes.PortalDataContainer{}
	err := get(fmt.Sprintf(portalsByName, mapId, portalName), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (m *mapInformation) GetPortals(mapId uint32) (*attributes.PortalDataContainer, error) {
	ar := &attributes.PortalDataContainer{}
	err := get(fmt.Sprintf(portalsResource, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
