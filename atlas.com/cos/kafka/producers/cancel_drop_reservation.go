package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type cancelDropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

var CancelDropReservation = func(l log.FieldLogger, ctx context.Context) *cancelDropReservation {
	return &cancelDropReservation{
		l:   l,
		ctx: ctx,
	}
}

type cancelDropReservation struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *cancelDropReservation) Emit(dropId uint32, characterId uint32) {
	event := &cancelDropReservationEvent{characterId, dropId}
	produceEvent(e.l, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND", createKey(int(dropId)), event)
}
