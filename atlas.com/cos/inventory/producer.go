package inventory

import (
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type inventoryModification struct {
	Mode          byte   `json:"mode"`
	ItemId        uint32 `json:"itemId"`
	InventoryType int8   `json:"inventoryType"`
	Quantity      uint32 `json:"quantity"`
	Position      int16  `json:"position"`
	OldPosition   int16  `json:"oldPosition"`
}

type inventoryModificationEvent struct {
	CharacterId   uint32                  `json:"characterId"`
	UpdateTick    bool                    `json:"updateTick"`
	Modifications []inventoryModification `json:"modifications"`
}

func emitInventoryModificationEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, updateTick bool, mode byte, itemId uint32, inventoryType int8, quantity uint32, position int16, oldPosition int16) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_INVENTORY_MODIFICATION")
	return func(characterId uint32, updateTick bool, mode byte, itemId uint32, inventoryType int8, quantity uint32, position int16, oldPosition int16) {
		event := &inventoryModificationEvent{
			CharacterId: characterId,
			UpdateTick:  updateTick,
			Modifications: []inventoryModification{
				{
					Mode:          mode,
					ItemId:        itemId,
					InventoryType: inventoryType,
					Quantity:      quantity,
					Position:      position,
					OldPosition:   oldPosition,
				},
			},
		}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type spawnDropCommand struct {
	WorldId      byte   `json:"worldId"`
	ChannelId    byte   `json:"channelId"`
	MapId        uint32 `json:"mapId"`
	ItemId       uint32 `json:"itemId"`
	EquipmentId  uint32 `json:"equipmentId"`
	Quantity     uint32 `json:"quantity"`
	Mesos        uint32 `json:"mesos"`
	DropType     byte   `json:"dropType"`
	X            int16  `json:"x"`
	Y            int16  `json:"y"`
	OwnerId      uint32 `json:"ownerId"`
	OwnerPartyId uint32 `json:"ownerPartyId"`
	DropperId    uint32 `json:"dropperId"`
	DropperX     int16  `json:"dropperX"`
	DropperY     int16  `json:"dropperY"`
	PlayerDrop   bool   `json:"playerDrop"`
	Mod          byte   `json:"mod"`
}

func emitSpawnItemDropCommand(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_SPAWN_CHARACTER_DROP_COMMAND")
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
		e := spawnDropCommand{
			WorldId:      worldId,
			ChannelId:    channelId,
			MapId:        mapId,
			ItemId:       itemId,
			EquipmentId:  0,
			Quantity:     quantity,
			Mesos:        mesos,
			DropType:     dropType,
			X:            x,
			Y:            y,
			OwnerId:      characterId,
			OwnerPartyId: characterPartyId,
			DropperId:    characterId,
			DropperX:     x,
			DropperY:     y,
			PlayerDrop:   true,
			Mod:          1,
		}
		producer(kafka.CreateKey(int(worldId)*1000+int(channelId)), e)
	}
}

func emitSpawnEquipDropCommand(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_SPAWN_CHARACTER_DROP_COMMAND")
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
		e := spawnDropCommand{
			WorldId:      worldId,
			ChannelId:    channelId,
			MapId:        mapId,
			ItemId:       itemId,
			EquipmentId:  equipmentId,
			Quantity:     0,
			Mesos:        0,
			DropType:     dropType,
			X:            x,
			Y:            y,
			OwnerId:      characterId,
			OwnerPartyId: characterPartyId,
			DropperId:    characterId,
			DropperX:     x,
			DropperY:     y,
			PlayerDrop:   true,
			Mod:          1,
		}
		producer(kafka.CreateKey(int(worldId)*1000+int(channelId)), e)
	}
}

type equipChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	Change      string `json:"change"`
}

func emitItemEquipped(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_EQUIP_CHANGED")
	return func(characterId uint32) {
		e := &equipChangedEvent{CharacterId: characterId, Change: "EQUIPPED"}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

func emitItemUnequipped(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_EQUIP_CHANGED")
	return func(characterId uint32) {
		e := &equipChangedEvent{CharacterId: characterId, Change: "UNEQUIPPED"}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
