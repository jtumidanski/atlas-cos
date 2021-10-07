package consumers

import (
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type gainItemCommand struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    int32  `json:"quantity"`
}

func GainItemCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &gainItemCommand{}
	}
}

func HandleGainItemCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*gainItemCommand); ok {
			if it, ok := inventory.GetInventoryType(event.ItemId); ok {
				if event.Quantity > 0 {
					err := item.GainItem(l, db, span)(event.CharacterId, it, event.ItemId, uint32(event.Quantity))
					if err != nil {
						l.WithError(err).Errorf("Unable to give character %d item %d.", event.CharacterId, event.ItemId)
					}
				} else {
					err := item.LoseItem(l, db, span)(event.CharacterId, it, event.ItemId, event.Quantity)
					if err != nil {
						l.WithError(err).Errorf("Unable to take item %d from character %d.", event.ItemId, event.CharacterId)
					}
				}
			} else {
				l.Errorf("Unable to locate inventory item %d belongs in.", event.ItemId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
