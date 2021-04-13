package consumers

import (
	"atlas-cos/character"
	"gorm.io/gorm"
	"log"
)

type changeMapCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
}

func ChangeMapCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &changeMapCommand{}
	}
}

func HandleChangeMapCommand(db *gorm.DB) EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*changeMapCommand); ok {
			character.Processor(l, db).ChangeMap(event.CharacterId, event.WorldId, event.ChannelId, event.MapId, event.PortalId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [ChangeMapCommand]")
		}
	}
}
