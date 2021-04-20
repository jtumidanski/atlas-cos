package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type characterStatUpdateEvent struct {
	CharacterId uint32   `json:"characterId"`
	Updates     []string `json:"updates"`
}

var CharacterStatUpdate = func(l log.FieldLogger, ctx context.Context) *characterStatUpdate {
	return &characterStatUpdate{
		l:   l,
		ctx: ctx,
	}
}

type characterStatUpdate struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *characterStatUpdate) Emit(characterId uint32, updates []string) {
	event := &characterStatUpdateEvent{characterId, updates}
	produceEvent(e.l, "TOPIC_CHARACTER_STAT_EVENT", createKey(int(characterId)), event)
}
