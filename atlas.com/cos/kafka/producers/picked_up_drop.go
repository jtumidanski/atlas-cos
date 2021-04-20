package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type pickedUpDropEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

var PickedUpDrop = func(l log.FieldLogger, ctx context.Context) *pickedUpDrop {
	return &pickedUpDrop{
		l:   l,
		ctx: ctx,
	}
}

type pickedUpDrop struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *pickedUpDrop) Emit(characterId uint32, dropId uint32) {
	event := &pickedUpDropEvent{characterId, dropId}
	produceEvent(e.l, "TOPIC_PICKUP_DROP_COMMAND", createKey(int(characterId)), event)
}
