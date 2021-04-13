package producers

import (
	"context"
	"log"
)

type mesoGainedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

var MesoGained = func(l *log.Logger, ctx context.Context) *mesoGained {
	return &mesoGained{
		l:   l,
		ctx: ctx,
	}
}

type mesoGained struct {
	l   *log.Logger
	ctx context.Context
}

func (e *mesoGained) Emit(characterId uint32, gain uint32) {
	event := &mesoGainedEvent{characterId, gain}
	produceEvent(e.l, "TOPIC_MESO_GAINED", createKey(int(characterId)), event)
}
