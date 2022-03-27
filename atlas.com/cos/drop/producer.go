package drop

import (
	"atlas-cos/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type adjustMesoCommand struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int32  `json:"amount"`
	Show        bool   `json:"show"`
}

func emitMesoAdjustment(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_ADJUST_MESO")
	return func(characterId uint32, amount int32) {
		event := &adjustMesoCommand{characterId, amount, true}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type cancelReservationCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func emitCancelReservation(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND")
	return func(dropId uint32, characterId uint32) {
		event := &cancelReservationCommand{characterId, dropId}
		producer(kafka.CreateKey(int(dropId)), event)
	}
}

type inventoryFullEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func emitInventoryFullEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_INVENTORY_FULL")
	return func(characterId uint32) {
		c := inventoryFullEvent{CharacterId: characterId}
		producer(kafka.CreateKey(int(characterId)), c)
	}
}

type pickedUpDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func emitPickedUpDropCommand(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, dropId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PICKUP_DROP_COMMAND")
	return func(characterId uint32, dropId uint32) {
		event := &pickedUpDropCommand{characterId, dropId}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type pickedUpItemEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func emitPickedUpItemEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, itemId uint32, quantity uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PICKED_UP_ITEM")
	return func(characterId uint32, itemId uint32, quantity uint32) {
		event := &pickedUpItemEvent{characterId, itemId, quantity}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}

type pickedUpNxEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func emitPickedUpNxEvent(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, gain uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PICKED_UP_NX")
	return func(characterId uint32, gain uint32) {
		event := &pickedUpNxEvent{characterId, gain}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}
