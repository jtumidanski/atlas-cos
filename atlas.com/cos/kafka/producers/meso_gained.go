package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type mesoGainedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

var MesoGained = func(l log.FieldLogger, ctx context.Context) *mesoGained {
	return &mesoGained{
		l:   l,
		ctx: ctx,
	}
}

type mesoGained struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *mesoGained) Emit(characterId uint32, gain int32) {
	event := &mesoGainedEvent{characterId, gain}
	produceEvent(e.l, "TOPIC_MESO_GAINED", createKey(int(characterId)), event)
}
