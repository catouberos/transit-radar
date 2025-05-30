package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/models"
	"github.com/catouberos/geoloc/protos"
	"github.com/catouberos/geoloc/rpc"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db *sql.DB

func main() {
	connection, exists := os.LookupEnv("DATABASE_CONNECTION")
	if !exists {
		log.Fatalln("Database connection has not been defined")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := models.New(conn)

	s := grpc.NewServer()

	app := base.InitApp(queries, s)

	protos.RegisterGeolocationServer(s, &rpc.GeolocationServer{
		App: app,
	})

	reflection.Register(s)

	if err := app.Start(); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
}
