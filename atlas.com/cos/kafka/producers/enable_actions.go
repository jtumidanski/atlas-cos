package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func EnableActions(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_ENABLE_ACTIONS")
	return func(characterId uint32) {
		event := &enableActionsEvent{characterId}
		producer(CreateKey(int(characterId)), event)
	}
}
