package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	monsterServicePrefix string = "/ms/mis/"
	monsterService              = baseRequest + monsterServicePrefix
	monstersResource            = monsterService + "monsters"
	monsterResource             = monstersResource + "/%d"
)

var Monster = func() *monster {
	return &monster{}
}

type monster struct {
}

func (m *monster) GetById(monsterId uint32) (*attributes.MonsterDataContainer, error) {
	ar := &attributes.MonsterDataContainer{}
	err := get(fmt.Sprintf(monsterResource, monsterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
