package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/catouberos/transit-radar/internal/base"
	"github.com/catouberos/transit-radar/internal/events"
	"github.com/catouberos/transit-radar/internal/queues"
	"github.com/catouberos/transit-radar/internal/server"
	"github.com/catouberos/transit-radar/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	done := make(chan bool)

	// set up database
	connection, exists := os.LookupEnv("DATABASE_CONNECTION")
	if !exists {
		log.Fatalln("Database connection has not been defined")
	}

	// setup main connection
	pool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		log.Panicln(err)
	}
	defer pool.Close()

	// setup connectrpc
	mux := server.NewRPCServer()

	// setup rabbitmq
	queue := queues.New("amqp://guest:guest@localhost:5672/")
	defer queue.Close()

	// initialise app
	app := base.NewApp(pool, migrations.Migrations, mux, queue)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		done <- true
	}()

	go func() {
		events.RegisterConsumer(app)

		done <- true
	}()

	// serve
	go func() {
		if err := app.Serve(); err != nil {
			log.Fatalln(err)
		}

		done <- true
	}()

	<-done
}
