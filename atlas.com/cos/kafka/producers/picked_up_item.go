package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type pickedUpItemEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

var PickedUpItem = func(l log.FieldLogger, ctx context.Context) *pickedUpItem {
	return &pickedUpItem{
		l:   l,
		ctx: ctx,
	}
}

type pickedUpItem struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *pickedUpItem) Emit(characterId uint32, itemId uint32, quantity uint32) {
	event := &pickedUpItemEvent{characterId, itemId, quantity}
	produceEvent(e.l, "TOPIC_PICKED_UP_ITEM", createKey(int(characterId)), event)
}
