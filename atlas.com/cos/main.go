package main

import (
	"atlas-cos/character"
	"atlas-cos/equipment"
	"atlas-cos/item"
	"atlas-cos/kafka/consumers"
	"atlas-cos/location"
	"atlas-cos/rest"
	"atlas-cos/retry"
	"atlas-cos/skill"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

func connectToDatabase(attempt int) (bool, interface{}, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:the@tcp(atlas-db:3306)/atlas-cos?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		return true, nil, err
	}
	return false, db, err
}

func main() {
	l := log.New()
	l.SetOutput(os.Stdout)
	if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := log.ParseLevel(val); err == nil {
			l.SetLevel(level)
		}
	}

	r, err := retry.RetryResponse(connectToDatabase, 10)
	if err != nil {
		panic("failed to connect database")
	}
	db := r.(*gorm.DB)

	// Migrate the schema
	character.Migration(db)
	equipment.Migration(db)
	item.Migration(db)
	location.Migration(db)
	skill.Migration(db)

	createEventConsumers(l, db)
	createRestService(l, db)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)
}

func createEventConsumers(l *log.Logger, db *gorm.DB) {
	cec := func(topicToken string, emptyEventCreator consumers.EmptyEventCreator, processor consumers.EventProcessor) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_ADJUST_HEALTH", consumers.AdjustHealthCommandCreator(), consumers.HandleAdjustHealthCommand(db))
	cec("TOPIC_ADJUST_MANA", consumers.AdjustManaCommandCreator(), consumers.HandleAdjustManaCommand(db))
	cec("TOPIC_ADJUST_MESO", consumers.AdjustMesoCommandCreator(), consumers.HandleAdjustMesoCommand(db))
	cec("TOPIC_ASSIGN_AP_COMMAND", consumers.AssignAPCommandCreator(), consumers.HandleAssignAPCommand(db))
	cec("TOPIC_ASSIGN_SP_COMMAND", consumers.AssignSPCommandCreator(), consumers.HandleAssignSPCommand(db))
	cec("TOPIC_CHANGE_MAP_COMMAND", consumers.ChangeMapCommandCreator(), consumers.HandleChangeMapCommand(db))
	cec("TOPIC_CHARACTER_EXPERIENCE_EVENT", consumers.GainExperienceEventCreator(), consumers.HandleGainExperienceEvent(db))
	cec("TOPIC_CHARACTER_LEVEL_EVENT", consumers.GainLevelEventCreator(), consumers.HandleGainLevelEvent(db))
	cec("TOPIC_CHARACTER_MOVEMENT", consumers.CharacterMovementEventCreator(), consumers.HandleCharacterMovementEvent(db))
	cec("TOPIC_CHARACTER_STATUS", consumers.CharacterStatusEventCreator(), consumers.HandleCharacterStatusEvent(db))
	cec("TOPIC_DROP_RESERVATION_EVENT", consumers.DropReservationEventCreator(), consumers.HandleDropReservationEvent(db))
	cec("TOPIC_MONSTER_KILLED_EVENT", consumers.MonsterKilledEventCreator(), consumers.HandleMonsterKilledEvent(db))
}

func createEventConsumer(l *log.Logger, topicToken string, emptyEventCreator consumers.EmptyEventCreator, processor consumers.EventProcessor) {
	h := func(logger log.FieldLogger, event interface{}) {
		processor(logger, event)
	}

	c := consumers.NewConsumer(l, context.Background(), h,
		consumers.SetGroupId("Character Orchestration Service"),
		consumers.SetTopicToken(topicToken),
		consumers.SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}

func createRestService(l *log.Logger, db *gorm.DB) {
	rs := rest.NewServer(l, db)
	go rs.Run()
}
