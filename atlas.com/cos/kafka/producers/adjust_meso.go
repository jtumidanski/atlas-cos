package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type adjustMesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int32  `json:"amount"`
	Show        bool   `json:"show"`
}

func AdjustMeso(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int32) {
	producer := ProduceEvent(l, span, "TOPIC_ADJUST_MESO")
	return func(characterId uint32, amount int32) {
		event := &adjustMesoEvent{characterId, amount, true}
		producer(CreateKey(int(characterId)), event)
	}
}
