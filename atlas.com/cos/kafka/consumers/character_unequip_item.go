package consumers

import (
	"atlas-cos/equipment"
	"atlas-cos/kafka/handler"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type characterUnequipItem struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func CharacterUnequipItemCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterUnequipItem{}
	}
}

func HandleCharacterUnequipItemCommand(db *gorm.DB) handler.EventHandler {
	return func(l log.FieldLogger, e interface{}) {
		if event, ok := e.(*characterUnequipItem); ok {
			l.Debugf("Begin event handling.")
			l.Debugf("CharacterId = %d, Source = %d, Destination = %d.", event.CharacterId, event.Source, event.CharacterId)
			e, err := equipment.GetEquippedItemForCharacterBySlot(l, db)(event.CharacterId, event.Source)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve item to equip for character %d in slot %d.", event.CharacterId, event.Source)
				return
			}
			equipment.UnequipItemForCharacter(l, db)(event.CharacterId, e.EquipmentId(), event.Source)
			l.Debugf("Complete event handling.")
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
