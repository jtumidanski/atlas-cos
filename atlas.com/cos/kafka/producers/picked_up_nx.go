package producers

import (
	"github.com/sirupsen/logrus"
)

type pickedUpNxEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func PickedUpNx(l logrus.FieldLogger) func(characterId uint32, gain uint32) {
	producer := ProduceEvent(l, "TOPIC_PICKED_UP_NX")
	return func(characterId uint32, gain uint32) {
		event := &pickedUpNxEvent{characterId, gain}
		producer(CreateKey(int(characterId)), event)
	}
}
