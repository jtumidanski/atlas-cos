package main

import (
	"atlas-cos/database"
	"atlas-cos/kafka/consumers"
	"atlas-cos/logger"
	"atlas-cos/rest"
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

	db := database.ConnectToDatabase(l)

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
