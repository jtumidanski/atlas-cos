package party

import (
	"atlas-cos/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	partyRegistryServicePrefix string = "/ms/party/"
	partyRegistryService              = requests.BaseRequest + partyRegistryServicePrefix
	charactersResource                = partyRegistryService + "characters"
	characterResource                 = charactersResource + "/%d"
	characterPartyResource            = characterResource + "/party"
)

type Request func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error)

func makeRequest(url string) Request {
	return func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(url, ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestByMemberId(memberId uint32) Request {
	return makeRequest(fmt.Sprintf(characterPartyResource, memberId))
}