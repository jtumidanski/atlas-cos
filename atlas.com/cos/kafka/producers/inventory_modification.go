package producers

import (
	"context"
	"log"
)

type inventoryModification struct {
	Mode          byte   `json:"mode"`
	ItemId        uint32 `json:"itemId"`
	InventoryType byte   `json:"inventoryType"`
	Quantity      uint32 `json:"quantity"`
	Position      int16  `json:"position"`
	OldPosition   int16  `json:"oldPosition"`
}

var InventoryModificationReservation = func(l *log.Logger, ctx context.Context) *inventoryModificationReservation {
	return &inventoryModificationReservation{
		l:   l,
		ctx: ctx,
	}
}

type characterInventoryModificationEvent struct {
	CharacterId   uint32                  `json:"characterId"`
	UpdateTick    bool                    `json:"updateTick"`
	Modifications []inventoryModification `json:"modifications"`
}

type inventoryModificationReservation struct {
	l   *log.Logger
	ctx context.Context
}

func (e *inventoryModificationReservation) Emit(characterId uint32, updateTick bool, mode byte, itemId uint32, inventoryType byte, quantity uint32, position int16) {
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
				OldPosition:   0,
			},
		},
	}
	produceEvent(e.l, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND", createKey(int(characterId)), event)
}
