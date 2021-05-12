package producers

import (
	"github.com/sirupsen/logrus"
)

type characterLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func CharacterLevel(l logrus.FieldLogger) func(characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_LEVEL_EVENT")
	return func(characterId uint32) {
		event := &characterLevelEvent{characterId}
		producer(CreateKey(int(characterId)), event)
	}
}
