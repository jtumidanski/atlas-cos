package character

import (
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type createdEvent struct {
	CharacterId uint32 `json:"characterId"`
	WorldId     byte   `json:"worldId"`
	Name        string `json:"name"`
}

func emitCreatedEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, worldId byte, name string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_CREATED_EVENT")
	return func(characterId uint32, worldId byte, name string) {
		event := &createdEvent{characterId, worldId, name}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type levelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func emitLevelEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_LEVEL_EVENT")
	return func(characterId uint32) {
		event := &levelEvent{characterId}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type statUpdateEvent struct {
	CharacterId uint32   `json:"characterId"`
	Updates     []string `json:"updates"`
}

func emitStatUpdateEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, updates []string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_STAT_EVENT")
	return func(characterId uint32, updates []string) {
		event := &statUpdateEvent{characterId, updates}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type enableActionsCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func emitEnableActionsCommand(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ENABLE_ACTIONS")
	return func(characterId uint32) {
		event := &enableActionsCommand{characterId}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func emitMapChangedEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHANGE_MAP_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
		event := &mapChangedEvent{worldId, channelId, mapId, portalId, characterId}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type mesoGainedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

func emitMesoGainedEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, gain int32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MESO_GAINED")
	return func(characterId uint32, gain int32) {
		event := &mesoGainedEvent{characterId, gain}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type skillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func emitSkillUpdateEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_SKILL_UPDATE_EVENT")
	return func(characterId uint32, skillId uint32, level uint32, masterLevel uint32, expiration int64) {
		event := &skillUpdateEvent{characterId, skillId, level, masterLevel, expiration}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}
