package information

import (
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	skillServicePrefix string = "/ms/sis/"
	skillService              = requests.BaseRequest + skillServicePrefix
	skillsResource            = skillService + "skills"
	skillResource             = skillsResource + "/%d"
)

func requestSkill(l logrus.FieldLogger, span opentracing.Span) func(skillId uint32) (*dataContainer, error) {
	return func(skillId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(skillResource, skillId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
