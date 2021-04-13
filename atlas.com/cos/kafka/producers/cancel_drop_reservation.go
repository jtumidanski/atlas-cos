package producers

import (
	"context"
	"log"
)

type cancelDropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

var CancelDropReservation = func(l *log.Logger, ctx context.Context) *cancelDropReservation {
	return &cancelDropReservation{
		l:   l,
		ctx: ctx,
	}
}

type cancelDropReservation struct {
	l   *log.Logger
	ctx context.Context
}

func (e *cancelDropReservation) Emit(dropId uint32, characterId uint32) {
	event := &cancelDropReservationEvent{characterId, dropId}
	produceEvent(e.l, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND", createKey(int(dropId)), event)
}
