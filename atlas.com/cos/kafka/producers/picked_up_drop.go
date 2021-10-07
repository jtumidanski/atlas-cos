package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type pickedUpDropEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func PickedUpDrop(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, dropId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_PICKUP_DROP_COMMAND")
	return func(characterId uint32, dropId uint32) {
		event := &pickedUpDropEvent{characterId, dropId}
		producer(CreateKey(int(characterId)), event)
	}
}
