package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/internal/events"
	"github.com/catouberos/geoloc/internal/queues"
	"github.com/catouberos/geoloc/internal/rpc"
	"github.com/catouberos/geoloc/protos"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	// setup grpc
	server := grpc.NewServer()
	defer server.Stop()

	// setup rabbitmq
	queue := queues.New("amqp://guest:guest@localhost:5672/")
	defer queue.Close()

	// initialise app
	app := base.NewApp(pool, embedMigrations, server, queue)

	protos.RegisterGeolocationServer(server, &rpc.GeolocationServer{
		App: app,
	})

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
