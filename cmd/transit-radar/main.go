package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/catouberos/transit-radar/internal/app"
	"github.com/catouberos/transit-radar/internal/config"
	"github.com/catouberos/transit-radar/internal/server"
	"github.com/catouberos/transit-radar/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(config.PostgresConnection, "postgresql://postgres:postgres@localhost:5432")
	viper.SetDefault(config.RedisConnection, "redis://localhost:6380")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// setup main connection
	pool, err := pgxpool.New(ctx, viper.GetString(config.PostgresConnection))
	if err != nil {
		log.Panicln(err)
	}
	defer pool.Close()

	// setup redis
	opt, err := redis.ParseURL(viper.GetString(config.RedisConnection))
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

	// initialise connectrpc server
	server := server.NewRPCServer(app, "[::]:5001")

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
