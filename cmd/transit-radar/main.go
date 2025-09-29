package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/catouberos/transit-radar/internal/app"
	"github.com/catouberos/transit-radar/internal/server"
	"github.com/catouberos/transit-radar/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

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

	// setup redis
	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		log.Panicln(err)
	}
	client := redis.NewClient(opt)
	defer client.Close()

	// initialise app
	app := app.NewApp(pool, migrations.Migrations, client)
	err = app.Init()
	if err != nil {
		log.Panicln(err)
	}

	// initialise grpc server
	server := server.NewRPCServer(app, "[::]:5000")

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

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
