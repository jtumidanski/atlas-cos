package producers

import "github.com/sirupsen/logrus"

type command struct {
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

func SpawnCharacterItemDrop(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
		e := command{
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
		produceEvent(l, "TOPIC_SPAWN_CHARACTER_DROP_COMMAND", createKey(int(worldId)*1000+int(channelId)), e)
	}
}

func SpawnCharacterEquipDrop(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, dropType byte, x int16, y int16, characterId uint32, characterPartyId uint32) {
		e := command{
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
		produceEvent(l, "TOPIC_SPAWN_CHARACTER_DROP_COMMAND", createKey(int(worldId)*1000+int(channelId)), e)
	}
}
