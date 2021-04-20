package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

var EnableActions = func(l log.FieldLogger, ctx context.Context) *enableActions {
	return &enableActions{
		l:   l,
		ctx: ctx,
	}
}

type enableActions struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *enableActions) Emit(characterId uint32) {
	event := &enableActionsEvent{characterId}
	produceEvent(e.l, "TOPIC_ENABLE_ACTIONS", createKey(int(characterId)), event)
}