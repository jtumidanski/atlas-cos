package producers

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type characterSkillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

var CharacterSkillUpdate = func(l log.FieldLogger, ctx context.Context) *characterSkillUpdate {
	return &characterSkillUpdate{
		l:   l,
		ctx: ctx,
	}
}

type characterSkillUpdate struct {
	l   log.FieldLogger
	ctx context.Context
}

func (c *characterSkillUpdate) Emit(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) {
	event := &characterSkillUpdateEvent{characterId, skillId, level, masterLevel, expiration}
	produceEvent(c.l, "TOPIC_CHARACTER_SKILL_UPDATE_EVENT", createKey(int(characterId)), event)
}
