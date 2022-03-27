package main

import (
	"atlas-cos/character"
	"atlas-cos/database"
	"atlas-cos/drop"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/kafka"
	"atlas-cos/location"
	"atlas-cos/logger"
	"atlas-cos/rest"
	"atlas-cos/seed"
	"atlas-cos/skill"
	"atlas-cos/tracing"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
import _ "net/http/pprof"

const serviceName = "atlas-cos"
const consumerGroupId = "Character Orchestration Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	db := database.Connect(l, database.SetMigrations(character.Migration, equipment.Migration, item.Migration, location.Migration, skill.Migration, inventory.Migration))

	kafka.CreateConsumers(l, ctx, wg,
		character.AdjustHealthCommandConsumer(db)(consumerGroupId),
		character.AdjustManaCommandConsumer(db)(consumerGroupId),
		character.AdjustMesoCommandConsumer(db)(consumerGroupId),
		character.AssignAPCommandConsumer(db)(consumerGroupId),
		character.AssignSPCommandConsumer(db)(consumerGroupId),
		character.ChangeMapCommandConsumer(db)(consumerGroupId),
		character.AdjustJobCommandConsumer(db)(consumerGroupId),
		character.DropItemCommandConsumer(db)(consumerGroupId),
		character.ResetAPCommandConsumer(db)(consumerGroupId),
		character.MoveItemCommandConsumer(db)(consumerGroupId),
		character.AdjustExperienceCommandConsumer(db)(consumerGroupId),
		character.GainLevelCommandConsumer(db)(consumerGroupId),
		character.MovementCommandConsumer(db)(consumerGroupId),
		character.StatusEventConsumer(db)(consumerGroupId),
		inventory.GainItemCommandConsumer(db)(consumerGroupId),
		drop.ReservationConsumer(db)(consumerGroupId),
		character.EquipItemConsumer(db)(consumerGroupId),
		character.UnEquipItemConsumer(db)(consumerGroupId),
		inventory.GainEquipmentConsumer(db)(consumerGroupId),
	)

	rest.CreateService(l, db, ctx, wg, "/ms/cos", character.InitResource, inventory.InitResource, seed.InitResource, location.InitResource, skill.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
