package consumers

import (
	"atlas-cos/monster"
	"gorm.io/gorm"
	"log"
)

type monsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         uint32        `json:"mapId"`
	UniqueId      uint32        `json:"uniqueId"`
	MonsterId     uint32        `json:"monsterId"`
	X             int16         `json:"x"`
	Y             int16         `json:"y"`
	KillerId      uint32        `json:"killerId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	Character uint32 `json:"character"`
	Damage    uint64  `json:"damage"`
}

func MonsterKilledEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &monsterKilledEvent{}
	}
}

func HandleMonsterKilledEvent(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*monsterKilledEvent); ok {
			if m, ok := monster.Processor(l, db).GetMonster(event.MonsterId); ok {
				var damageEntries = make([]*monster.DamageEntry, 0)
				for _, entry := range event.DamageEntries {
					damageEntries = append(damageEntries, monster.NewDamageEntry(entry.Character, entry.Damage))
				}
				monster.Processor(l, db).DistributeExperience(event.WorldId, event.ChannelId, event.MapId, m, damageEntries)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [MonsterKilledEvent]")
		}
	}
}
