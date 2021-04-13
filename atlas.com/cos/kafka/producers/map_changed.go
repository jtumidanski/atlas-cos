package producers

import (
	"context"
	"log"
)

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

var MapChanged = func(l *log.Logger, ctx context.Context) *mapChanged {
	return &mapChanged{
		l:   l,
		ctx: ctx,
	}
}

type mapChanged struct {
	l   *log.Logger
	ctx context.Context
}

func (e *mapChanged) Emit(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	event := &mapChangedEvent{worldId, channelId, mapId, portalId, characterId}
	produceEvent(e.l, "TOPIC_CHANGE_MAP_EVENT", createKey(int(characterId)), event)
}
