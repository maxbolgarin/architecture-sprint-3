package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Println("starting telemetry microservice")

	db, err := NewDB(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if err := StartKafka(ctx, db); err != nil {
		log.Fatalln(err)
	}

	StartApp(ctx, db)

	<-ctx.Done()
}
