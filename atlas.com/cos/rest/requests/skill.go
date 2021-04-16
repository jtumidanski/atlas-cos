package requests

import (
	"atlas-cos/rest/attributes"
	"fmt"
)

const (
	skillServicePrefix string = "/ms/sis/"
	skillService              = baseRequest + skillServicePrefix
	skillsResource            = skillService + "skills"
	skillResource             = skillsResource + "/%d"
)

var Skill = func() *skill {
	return &skill{}
}

type skill struct {
}

func (m *skill) GetById(skillId uint32) (*attributes.SkillDataContainer, error) {
	ar := &attributes.SkillDataContainer{}
	err := get(fmt.Sprintf(skillResource, skillId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
