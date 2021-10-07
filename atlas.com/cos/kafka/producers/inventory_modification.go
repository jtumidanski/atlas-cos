package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type inventoryModification struct {
	Mode          byte   `json:"mode"`
	ItemId        uint32 `json:"itemId"`
	InventoryType int8   `json:"inventoryType"`
	Quantity      uint32 `json:"quantity"`
	Position      int16  `json:"position"`
	OldPosition   int16  `json:"oldPosition"`
}

type characterInventoryModificationEvent struct {
	CharacterId   uint32                  `json:"characterId"`
	UpdateTick    bool                    `json:"updateTick"`
	Modifications []inventoryModification `json:"modifications"`
}

func InventoryModificationReservation(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, updateTick bool, mode byte, itemId uint32, inventoryType int8, quantity uint32, position int16, oldPosition int16) {
	producer := ProduceEvent(l, span, "TOPIC_INVENTORY_MODIFICATION")
	return func(characterId uint32, updateTick bool, mode byte, itemId uint32, inventoryType int8, quantity uint32, position int16, oldPosition int16) {
		event := &characterInventoryModificationEvent{
			CharacterId: characterId,
			UpdateTick:  updateTick,
			Modifications: []inventoryModification{
				{
					Mode:          mode,
					ItemId:        itemId,
					InventoryType: inventoryType,
					Quantity:      quantity,
					Position:      position,
					OldPosition:   oldPosition,
				},
			},
		}
		producer(CreateKey(int(characterId)), event)
	}
}
