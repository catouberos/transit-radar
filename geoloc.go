package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/protos"
	"github.com/catouberos/geoloc/rpc"
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

	db, err := sqlx.Connect("pgx", connection)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	s := grpc.NewServer()

	app := base.InitApp(db.DB, s)

	protos.RegisterGeolocationServer(s, &rpc.GeolocationServer{
		App: app,
	})

	reflection.Register(s)

	if err := app.Start(); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
}
