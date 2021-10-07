package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type pickedUpNxEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func PickedUpNx(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, gain uint32) {
	producer := ProduceEvent(l, span, "TOPIC_PICKED_UP_NX")
	return func(characterId uint32, gain uint32) {
		event := &pickedUpNxEvent{characterId, gain}
		producer(CreateKey(int(characterId)), event)
	}
}
