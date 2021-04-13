package producers

import (
	"context"
	"log"
)

type pickedUpNxEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

var PickedUpNx = func(l *log.Logger, ctx context.Context) *pickedUpNx {
	return &pickedUpNx{
		l:   l,
		ctx: ctx,
	}
}

type pickedUpNx struct {
	l   *log.Logger
	ctx context.Context
}

func (e *pickedUpNx) Emit(characterId uint32, gain uint32) {
	event := &pickedUpNxEvent{characterId, gain}
	produceEvent(e.l, "TOPIC_PICKED_UP_NX", createKey(int(characterId)), event)
}
