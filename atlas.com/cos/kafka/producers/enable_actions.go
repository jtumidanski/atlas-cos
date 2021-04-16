package producers

import (
	"context"
	"log"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

var EnableActions = func(l *log.Logger, ctx context.Context) *enableActions {
	return &enableActions{
		l:   l,
		ctx: ctx,
	}
}

type enableActions struct {
	l   *log.Logger
	ctx context.Context
}

func (e *enableActions) Emit(characterId uint32) {
	event := &enableActionsEvent{characterId}
	produceEvent(e.l, "TOPIC_ENABLE_ACTIONS", createKey(int(characterId)), event)
}