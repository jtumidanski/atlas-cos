package main

import (
	"atlas-cos/character"
	"atlas-cos/database"
	"atlas-cos/equipment"
	"atlas-cos/inventory"
	"atlas-cos/item"
	"atlas-cos/kafka/consumers"
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

	consumers.CreateEventConsumers(l, db, ctx, wg)

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
