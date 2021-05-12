package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type changeMapCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
}

func ChangeMapCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &changeMapCommand{}
	}
}

func HandleChangeMapCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*changeMapCommand); ok {
			l.Debugf("Begin event handling.")
			character.ChangeMap(l, db)(event.CharacterId, event.WorldId, event.ChannelId, event.MapId, event.PortalId)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
