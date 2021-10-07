package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func MapChanged(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_CHANGE_MAP_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
		event := &mapChangedEvent{worldId, channelId, mapId, portalId, characterId}
		producer(CreateKey(int(characterId)), event)
	}
}
