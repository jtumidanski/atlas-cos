package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type adjustMesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint32 `json:"amount"`
	Show        bool   `json:"show"`
}

var AdjustMeso = func(l log.FieldLogger, ctx context.Context) *adjustMeso {
	return &adjustMeso{
		l:   l,
		ctx: ctx,
	}
}

type adjustMeso struct {
	l   log.FieldLogger
	ctx context.Context
}

func (e *adjustMeso) Emit(characterId uint32, amount uint32) {
	event := &adjustMesoEvent{characterId, amount, true}
	produceEvent(e.l, "TOPIC_ADJUST_MESO", createKey(int(characterId)), event)
}
