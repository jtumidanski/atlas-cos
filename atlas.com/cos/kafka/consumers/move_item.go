package consumers

import (
	"atlas-cos/character"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type moveItemCommand struct {
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Destination   int16  `json:"destination"`
}

func MoveItemCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &moveItemCommand{}
	}
}

func HandleMoveItemCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*moveItemCommand); ok {
			_ = character.MoveItem(l, db)(event.CharacterId, event.InventoryType, event.Source, event.Destination)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}