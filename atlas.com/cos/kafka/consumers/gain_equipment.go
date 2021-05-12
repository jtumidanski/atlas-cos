package consumers

import (
	"atlas-cos/equipment"
	"atlas-cos/rest/requests"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type gainEquipmentCommand struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
}

func GainEquipmentCommandCreator() EmptyEventCreator {
	return func() interface{} {
		return &gainEquipmentCommand{}
	}
}

func HandleGainEquipmentCommand(db *gorm.DB) EventProcessor {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*gainEquipmentCommand); ok {
			ro, err := requests.EquipmentRegistry().Create(event.ItemId)
			if err != nil {
				l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
				return
			}
			eid, err := strconv.Atoi(ro.Data.Id)
			if err != nil {
				l.Errorf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.")
				return
			}

			err = equipment.GainItem(l, db)(event.CharacterId, event.ItemId, uint32(eid))
			if err != nil {
				l.WithError(err).Errorf("Unable to give character %d item %d.", event.CharacterId, event.ItemId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
