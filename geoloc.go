package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/internal/rpc"
	"github.com/catouberos/geoloc/models"
	"github.com/catouberos/geoloc/protos"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// set up database
	connection, exists := os.LookupEnv("DATABASE_CONNECTION")
	if !exists {
		log.Fatalln("Database connection has not been defined")
	}

	// runs migration if present
	db, err := sql.Open("pgx", connection)
	if err != nil {
		log.Fatalln(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalln(err)
	}

	log.Println("Running migrations...")

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalln(err)
	}

	log.Println("Migration completed!")

	db.Close()

	// setup main connection
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	// setup models
	queries := models.New(conn)

	// setup grpc
	s := grpc.NewServer()

	// initialise app
	app := base.InitApp(queries, s)

	protos.RegisterGeolocationServer(s, &rpc.GeolocationServer{
		App: app,
	})

	// grpc: enable reflection
	reflection.Register(s)

	// serve
	if err := app.Start(); err != nil {
		log.Fatalln(err)
	}
}
