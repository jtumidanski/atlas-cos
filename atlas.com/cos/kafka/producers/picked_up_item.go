package producers

import (
	"github.com/sirupsen/logrus"
)

type pickedUpItemEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func PickedUpItem(l logrus.FieldLogger) func(characterId uint32, itemId uint32, quantity uint32) {
	producer := ProduceEvent(l, "TOPIC_PICKED_UP_ITEM")
	return func(characterId uint32, itemId uint32, quantity uint32) {
		event := &pickedUpItemEvent{characterId, itemId, quantity}
		producer(CreateKey(int(characterId)), event)
	}
}
