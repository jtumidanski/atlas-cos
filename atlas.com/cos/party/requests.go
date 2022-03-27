package party

import (
	"atlas-cos/rest/requests"
	"fmt"
)

const (
	partyRegistryServicePrefix string = "/ms/party/"
	partyRegistryService              = requests.BaseRequest + partyRegistryServicePrefix
	charactersResource                = partyRegistryService + "characters"
	characterResource                 = charactersResource + "/%d"
	characterPartyResource            = characterResource + "/party"
)

func requestByMemberId(memberId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(characterPartyResource, memberId))
}
