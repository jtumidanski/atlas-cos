package information

import (
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	skillServicePrefix string = "/ms/sis/"
	skillService              = requests.BaseRequest + skillServicePrefix
	skillsResource            = skillService + "skills"
	skillResource             = skillsResource + "/%d"
)

func requestSkill(l logrus.FieldLogger) func(skillId uint32) (*dataContainer, error) {
	return func(skillId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(skillResource, skillId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
