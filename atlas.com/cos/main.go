package main

import (
	"atlas-cos/database"
	"atlas-cos/kafka/consumers"
	"atlas-cos/logger"
	"atlas-cos/rest"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import _ "net/http/pprof"

func main() {
	l := logger.CreateLogger()

	go func() {
		l.Fatal(http.ListenAndServe(":6060", http.DefaultServeMux))
	}()

	db := database.ConnectToDatabase(l)

	consumers.CreateEventConsumers(l, db)

	rest.CreateRestService(l, db)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)
}
