package producers

import (
	"github.com/sirupsen/logrus"
)

type characterStatUpdateEvent struct {
	CharacterId uint32   `json:"characterId"`
	Updates     []string `json:"updates"`
}

func CharacterStatUpdate(l logrus.FieldLogger) func(characterId uint32, updates []string) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_STAT_EVENT")
	return func(characterId uint32, updates []string) {
		event := &characterStatUpdateEvent{characterId, updates}
		producer(CreateKey(int(characterId)), event)
	}
}
