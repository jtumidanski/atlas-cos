package inventory

import (
	"atlas-cos/equipment/statistics"
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	consumerNameGainItem = "gain_item_command"
	topicTokenGainItem   = "TOPIC_CHARACTER_GAIN_ITEM"
)

func GainItemCommandConsumer(db *gorm.DB) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[gainItemCommand](consumerNameGainItem, topicTokenGainItem, groupId, handleGainItemCommand(db))
	}
}

type gainItemCommand struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    int32  `json:"quantity"`
}

func handleGainItemCommand(db *gorm.DB) kafka.HandlerFunc[gainItemCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event gainItemCommand) {
		if it, ok := GetInventoryType(event.ItemId); ok {
			if it == TypeValueEquip {
				eid, err := statistics.Create(l, span)(event.ItemId)
				if err != nil {
					l.WithError(err).Errorf("Unable to create equipment %d for character %d.", event.ItemId, event.CharacterId)
					return
				}

				err = GainEquipment(l, db, span)(event.CharacterId, event.ItemId, eid)
				if err != nil {
					l.WithError(err).Errorf("Unable to give character %d item %d.", event.CharacterId, event.ItemId)
				}
			} else {
				if event.Quantity > 0 {
					err := GainItem(l, db, span)(event.CharacterId, it, event.ItemId, uint32(event.Quantity))
					if err != nil {
						l.WithError(err).Errorf("Unable to give character %d item %d.", event.CharacterId, event.ItemId)
					}
				} else {
					err := LoseItem(l, db, span)(event.CharacterId, it, event.ItemId, event.Quantity)
					if err != nil {
						l.WithError(err).Errorf("Unable to take item %d from character %d.", event.ItemId, event.CharacterId)
					}
				}
			}
		} else {
			l.Errorf("Unable to locate inventory item %d belongs in.", event.ItemId)
		}
	}
}
