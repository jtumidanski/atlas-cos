package producers

import (
	"github.com/sirupsen/logrus"
)

type characterSkillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func CharacterSkillUpdate(l logrus.FieldLogger) func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_SKILL_UPDATE_EVENT")
	return func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) {
		event := &characterSkillUpdateEvent{characterId, skillId, level, masterLevel, expiration}
		producer(CreateKey(int(characterId)), event)
	}
}
