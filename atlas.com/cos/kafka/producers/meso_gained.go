package producers

import (
	"github.com/sirupsen/logrus"
)

type mesoGainedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

func MesoGained(l logrus.FieldLogger) func(characterId uint32, gain int32) {
	producer := ProduceEvent(l, "TOPIC_MESO_GAINED")
	return func(characterId uint32, gain int32) {
		event := &mesoGainedEvent{characterId, gain}
		producer(CreateKey(int(characterId)), event)
	}
}
