package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type inventoryFullCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func InventoryFull(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_INVENTORY_FULL")
	return func(characterId uint32) {
		c := inventoryFullCommand{CharacterId: characterId}
		producer(CreateKey(int(characterId)), c)
	}
}
