package information

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	skillServicePrefix string = "/ms/sis/"
	skillService              = requests.BaseRequest + skillServicePrefix
	skillsResource            = skillService + "skills"
	skillResource             = skillsResource + "/%d"
)

func requestSkill(skillId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(skillResource, skillId))
}
