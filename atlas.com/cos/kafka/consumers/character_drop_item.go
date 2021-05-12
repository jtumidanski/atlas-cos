package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type characterDropItem struct {
	WorldId       byte   `json:"worldId"`
	ChannelId     byte   `json:"channelId"`
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Quantity      int16  `json:"quantity"`
}

func CharacterDropItemCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterDropItem{}
	}
}

func HandleCharacterDropItemCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*characterDropItem); ok {
			l.Debugf("Begin event handling.")
			l.Debugf("Request to drop %d item in slot %d for character %d.", event.Quantity, event.Source, event.CharacterId)
			_ = character.DropItem(l, db)(event.WorldId, event.ChannelId, event.CharacterId, event.InventoryType, event.Source, event.Quantity)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
