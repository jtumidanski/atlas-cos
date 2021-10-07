package consumers

import (
	"atlas-cos/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

const (
	AdjustHealthCommand      = "adjust_health_command"
	AdjustManaCommand        = "adjust_mana_command"
	AdjustMesoCommand        = "adjust_meso_command"
	AssignAPCommand          = "assign_ap_command"
	AssignSPCommand          = "assign_sp_command"
	ChangeMapCommand         = "change_map_command"
	CharacterExperienceEvent = "character_experience_event"
	CharacterLevelEvent      = "character_level_event"
	CharacterMovementEvent   = "character_movement_event"
	CharacterStatusEvent     = "character_status_event"
	DropReservationEvent     = "drop_reservation_event"
	EquipItemCommand         = "equip_item_command"
	UnequipItemCommand       = "unequip_item_command"
	GainEquipmentCommand     = "gain_equipment_command"
	GainItemCommand          = "gain_item_command"
	AdjustJobCommand         = "adjust_job_command"
	ResetAPCommand           = "reset_ap_command"
	DropItemCommand          = "drop_item_command"
	MoveItemCommand          = "move_item_command"
)

func CreateEventConsumers(l *logrus.Logger, db *gorm.DB, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_ADJUST_HEALTH", AdjustHealthCommand, AdjustHealthCommandCreator(), HandleAdjustHealthCommand(db))
	cec("TOPIC_ADJUST_MANA", AdjustManaCommand, AdjustManaCommandCreator(), HandleAdjustManaCommand(db))
	cec("TOPIC_ADJUST_MESO", AdjustMesoCommand, AdjustMesoCommandCreator(), HandleAdjustMesoCommand(db))
	cec("TOPIC_ASSIGN_AP_COMMAND", AssignAPCommand, AssignAPCommandCreator(), HandleAssignAPCommand(db))
	cec("TOPIC_ASSIGN_SP_COMMAND", AssignSPCommand, AssignSPCommandCreator(), HandleAssignSPCommand(db))
	cec("TOPIC_CHANGE_MAP_COMMAND", ChangeMapCommand, ChangeMapCommandCreator(), HandleChangeMapCommand(db))
	cec("TOPIC_CHARACTER_EXPERIENCE_EVENT", CharacterExperienceEvent, GainExperienceEventCreator(), HandleGainExperienceEvent(db))
	cec("TOPIC_CHARACTER_LEVEL_EVENT", CharacterLevelEvent, GainLevelEventCreator(), HandleGainLevelEvent(db))
	cec("TOPIC_CHARACTER_MOVEMENT", CharacterMovementEvent, CharacterMovementEventCreator(), HandleCharacterMovementEvent(db))
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusEvent, CharacterStatusEventCreator(), HandleCharacterStatusEvent(db))
	cec("TOPIC_DROP_RESERVATION_EVENT", DropReservationEvent, DropReservationEventCreator(), HandleDropReservationEvent(db))
	cec("TOPIC_EQUIP_ITEM", EquipItemCommand, CharacterEquipItemCommandCreator(), HandleCharacterEquipItemCommand(db))
	cec("TOPIC_UNEQUIP_ITEM", UnequipItemCommand, CharacterUnequipItemCommandCreator(), HandleCharacterUnequipItemCommand(db))
	cec("TOPIC_CHARACTER_GAIN_EQUIPMENT", GainEquipmentCommand, GainEquipmentCommandCreator(), HandleGainEquipmentCommand(db))
	cec("TOPIC_CHARACTER_GAIN_ITEM", GainItemCommand, GainItemCommandCreator(), HandleGainItemCommand(db))
	cec("TOPIC_CHARACTER_ADJUST_JOB", AdjustJobCommand, AdjustJobCommandCreator(), HandleAdjustJobCommand(db))
	cec("TOPIC_CHARACTER_RESET_AP", ResetAPCommand, ResetAPCommandCreator(), HandleResetAPCommand(db))
	cec("TOPIC_CHARACTER_DROP_ITEM", DropItemCommand, CharacterDropItemCommandCreator(), HandleCharacterDropItemCommand(db))
	cec("TOPIC_MOVE_ITEM", MoveItemCommand, MoveItemCommandCreator(), HandleMoveItemCommand(db))
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "Character Orchestration Service", emptyEventCreator, processor)
}
