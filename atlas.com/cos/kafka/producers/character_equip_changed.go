package producers

import log "github.com/sirupsen/logrus"

type characterEquipChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func CharacterEquippedItem(l log.FieldLogger) func(characterId uint32) {
	return func(characterId uint32) {
		e := &characterEquipChangedEvent{CharacterId: characterId, Change: "EQUIPPED"}
		produceEvent(l, "TOPIC_CHARACTER_EQUIP_CHANGED", createKey(int(characterId)), e)
	}
}

func CharacterUnEquippedItem(l log.FieldLogger) func(characterId uint32) {
	return func(characterId uint32) {
		e := &characterEquipChangedEvent{CharacterId: characterId, Change: "UNEQUIPPED"}
		produceEvent(l, "TOPIC_CHARACTER_EQUIP_CHANGED", createKey(int(characterId)), e)
	}
}
