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
	"atlas-cos/skill"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
import _ "net/http/pprof"

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	db := database.Connect(l, database.SetMigrations(character.Migration, equipment.Migration, item.Migration, location.Migration, skill.Migration, inventory.Migration))

	consumers.CreateEventConsumers(l, db, ctx, wg)

	rest.CreateRestService(l, db, ctx, wg)

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
