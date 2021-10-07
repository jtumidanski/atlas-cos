package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type characterLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func CharacterLevel(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_CHARACTER_LEVEL_EVENT")
	return func(characterId uint32) {
		event := &characterLevelEvent{characterId}
		producer(CreateKey(int(characterId)), event)
	}
}
