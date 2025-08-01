package base

import (
	"net"

	"github.com/catouberos/geoloc/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type App struct {
	Queries *models.Queries
	GRPC    *grpc.Server
	AMQP    *amqp.Channel
}

func InitApp(queries *models.Queries, grpc *grpc.Server, amqp *amqp.Channel) *App {
	app := &App{
		Queries: queries,
		GRPC:    grpc,
		AMQP:    amqp,
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
