package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService = baseRequest + dropRegistryServicePrefix
	dropResource        = dropRegistryService + "drops/%d"
)

var DropRegistry = func() *dropRegistry {
	return &dropRegistry{}
}

type dropRegistry struct {
}

func (d *dropRegistry) GetDropById(dropId uint32) (*attributes.DropDataContainer, error) {
	ar := &attributes.DropDataContainer{}
	err := get(fmt.Sprintf(dropResource, dropId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}