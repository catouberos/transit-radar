package base

import (
	"database/sql"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	DB   *sql.DB
	GRPC *grpc.Server
}

func InitApp(db *sql.DB, grpc *grpc.Server) *App {
	app := &App{
		DB:   db,
		GRPC: grpc,
	}

	return app
}

func (app *App) Start() error {

	lis, err := net.Listen("tcp", ":5005")
	if err != nil {
		return err
	}

	if err := app.GRPC.Serve(lis); err != nil {
		return err
	}

	return nil
}
