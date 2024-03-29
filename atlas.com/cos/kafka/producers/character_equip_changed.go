package producers

import "github.com/sirupsen/logrus"

type characterEquipChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func CharacterEquippedItem(l logrus.FieldLogger) func(characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_EQUIP_CHANGED")
	return func(characterId uint32) {
		e := &characterEquipChangedEvent{CharacterId: characterId, Change: "EQUIPPED"}
		producer(CreateKey(int(characterId)), e)
	}
}

func CharacterUnEquippedItem(l logrus.FieldLogger) func(characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_EQUIP_CHANGED")
	return func(characterId uint32) {
		e := &characterEquipChangedEvent{CharacterId: characterId, Change: "UNEQUIPPED"}
		producer(CreateKey(int(characterId)), e)
	}
}
