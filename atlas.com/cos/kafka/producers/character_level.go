package producers

import (
	"context"
	"log"
)

type characterLevelEvent struct {
	CharacterId uint32   `json:"characterId"`
}

var CharacterLevel = func(l *log.Logger, ctx context.Context) *characterLevel {
	return &characterLevel{
		l:   l,
		ctx: ctx,
	}
}

type characterLevel struct {
	l   *log.Logger
	ctx context.Context
}

func (e *characterLevel) Emit(characterId uint32) {
	event := &characterLevelEvent{characterId}
	produceEvent(e.l, "TOPIC_CHARACTER_LEVEL_EVENT", createKey(int(characterId)), event)
}