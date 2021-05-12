package consumers

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateEventConsumers(l *logrus.Logger, db *gorm.DB) {
	cec := func(topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_ADJUST_HEALTH", AdjustHealthCommandCreator(), HandleAdjustHealthCommand(db))
	cec("TOPIC_ADJUST_MANA", AdjustManaCommandCreator(), HandleAdjustManaCommand(db))
	cec("TOPIC_ADJUST_MESO", AdjustMesoCommandCreator(), HandleAdjustMesoCommand(db))
	cec("TOPIC_ASSIGN_AP_COMMAND", AssignAPCommandCreator(), HandleAssignAPCommand(db))
	cec("TOPIC_ASSIGN_SP_COMMAND", AssignSPCommandCreator(), HandleAssignSPCommand(db))
	cec("TOPIC_CHANGE_MAP_COMMAND", ChangeMapCommandCreator(), HandleChangeMapCommand(db))
	cec("TOPIC_CHARACTER_EXPERIENCE_EVENT", GainExperienceEventCreator(), HandleGainExperienceEvent(db))
	cec("TOPIC_CHARACTER_LEVEL_EVENT", GainLevelEventCreator(), HandleGainLevelEvent(db))
	cec("TOPIC_CHARACTER_MOVEMENT", CharacterMovementEventCreator(), HandleCharacterMovementEvent(db))
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusEventCreator(), HandleCharacterStatusEvent(db))
	cec("TOPIC_DROP_RESERVATION_EVENT", DropReservationEventCreator(), HandleDropReservationEvent(db))
	cec("TOPIC_EQUIP_ITEM", CharacterEquipItemCommandCreator(), HandleCharacterEquipItemCommand(db))
	cec("TOPIC_UNEQUIP_ITEM", CharacterUnequipItemCommandCreator(), HandleCharacterUnequipItemCommand(db))
	cec("TOPIC_CHARACTER_GAIN_EQUIPMENT", GainEquipmentCommandCreator(), HandleGainEquipmentCommand(db))
	cec("TOPIC_CHARACTER_GAIN_ITEM", GainItemCommandCreator(), HandleGainItemCommand(db))
	cec("TOPIC_CHARACTER_ADJUST_JOB", AdjustJobCommandCreator(), HandleAdjustJobCommand(db))
	cec("TOPIC_CHARACTER_RESET_AP", ResetAPCommandCreator(), HandleResetAPCommand(db))
	cec("TOPIC_CHARACTER_DROP_ITEM", CharacterDropItemCommandCreator(), HandleCharacterDropItemCommand(db))
}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
	h := func(logger logrus.FieldLogger, event interface{}) {
		processor(logger, event)
	}

	c := NewConsumer(l, context.Background(), h,
		SetGroupId("Character Orchestration Service"),
		SetTopicToken(topicToken),
		SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}
