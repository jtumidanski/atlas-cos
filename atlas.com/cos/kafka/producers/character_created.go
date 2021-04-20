package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type characterCreatedEvent struct {
	CharacterId uint32 `json:"characterId"`
	WorldId     byte   `json:"worldId"`
	Name        string `json:"name"`
}

var CharacterCreated = func(l log.FieldLogger, ctx context.Context) *characterCreated {
	return &characterCreated{
		l:   l,
		ctx: ctx,
	}
}

type characterCreated struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *characterCreated) Emit(characterId uint32, worldId byte, name string) {
	event := &characterCreatedEvent{characterId, worldId, name}
	produceEvent(e.l, "TOPIC_CHARACTER_CREATED_EVENT", createKey(int(characterId)), event)
}
