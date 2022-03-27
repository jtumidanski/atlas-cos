package equipment

import (
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

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
