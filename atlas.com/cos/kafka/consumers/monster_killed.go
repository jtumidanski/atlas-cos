package consumers

import (
	"gorm.io/gorm"
	"log"
)

type monsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         int           `json:"mapId"`
	UniqueId      int           `json:"uniqueId"`
	MonsterId     int           `json:"monsterId"`
	X             int           `json:"x"`
	Y             int           `json:"y"`
	KillerId      int           `json:"killerId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	Character int   `json:"character"`
	Damage    int64 `json:"damage"`
}

func MonsterKilledEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &monsterKilledEvent{}
	}
}

func HandleMonsterKilledEvent(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*monsterKilledEvent); ok {

		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [MonsterKilledEvent]")
		}
	}
}
