package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type pickedUpItemEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func PickedUpItem(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, itemId uint32, quantity uint32) {
	producer := ProduceEvent(l, span, "TOPIC_PICKED_UP_ITEM")
	return func(characterId uint32, itemId uint32, quantity uint32) {
		event := &pickedUpItemEvent{characterId, itemId, quantity}
		producer(CreateKey(int(characterId)), event)
	}
}
