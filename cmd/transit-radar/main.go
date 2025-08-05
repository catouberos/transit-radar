package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/catouberos/transit-radar/internal/base"
	"github.com/catouberos/transit-radar/internal/queue"
	"github.com/catouberos/transit-radar/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wagslane/go-rabbitmq"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// set up database
	connection, exists := os.LookupEnv("DATABASE_CONNECTION")
	if !exists {
		log.Fatalln("Database connection has not been defined")
	}

	// setup main connection
	pool, err := pgxpool.New(ctx, connection)
	if err != nil {
		log.Panicln(err)
	}
	defer pool.Close()

	// setup rabbitmq
	rmq, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Panicln(err)
	}
	defer rmq.Close()

	// initialise app
	app := base.NewApp(pool, migrations)
	err = app.Init()
	if err != nil {
		log.Panicln(err)
	}

	// initialise grpc server
	server := server.NewRPCServer(app, "[::]:5000")

	// initialist rabbitmq
	handler := queue.NewConsumerHandler(rmq, app)
	defer handler.Stop()

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		cancel()
	}()

	go func() {
		handler.Start()

		cancel()
	}()

	// serve
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalln(err)
		}

		cancel()
	}()

	<-ctx.Done()
}
