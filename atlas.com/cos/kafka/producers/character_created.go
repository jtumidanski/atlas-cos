package producers

import (
	"github.com/sirupsen/logrus"
)

type characterCreatedEvent struct {
	CharacterId uint32 `json:"characterId"`
	WorldId     byte   `json:"worldId"`
	Name        string `json:"name"`
}

func CharacterCreated(l logrus.FieldLogger) func(characterId uint32, worldId byte, name string) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_CREATED_EVENT")
	return func(characterId uint32, worldId byte, name string) {
		event := &characterCreatedEvent{characterId, worldId, name}
		producer(CreateKey(int(characterId)), event)
	}
}
