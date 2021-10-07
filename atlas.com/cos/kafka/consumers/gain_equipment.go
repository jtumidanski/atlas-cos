package consumers

import (
	"atlas-cos/equipment"
	"atlas-cos/equipment/statistics"
	"atlas-cos/kafka/handler"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type gainEquipmentCommand struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
}

func GainEquipmentCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &gainEquipmentCommand{}
	}
}

func HandleGainEquipmentCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*gainEquipmentCommand); ok {
			eid, err := statistics.Create(l, span)(event.ItemId)
			if err != nil {
				l.WithError(err).Errorf("Unable to create equipment %d for character %d.", event.ItemId, event.CharacterId)
				return
			}

			err = equipment.GainItem(l, db, span)(event.CharacterId, event.ItemId, eid)
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", event.CharacterId, event.ItemId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
