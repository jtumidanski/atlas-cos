package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type mesoGainedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

func MesoGained(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, gain int32) {
	producer := ProduceEvent(l, span, "TOPIC_MESO_GAINED")
	return func(characterId uint32, gain int32) {
		event := &mesoGainedEvent{characterId, gain}
		producer(CreateKey(int(characterId)), event)
	}
}
