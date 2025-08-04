package events

import (
	"context"
	"time"

	"github.com/catouberos/geoloc/base"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RegisterConsumer(app *base.App) {
	for {
		ctx, cancel := context.WithCancel(context.Background())

		for {
			if app.Queue().IsReady() {
				break
			}

			<-time.After(2 * time.Second)
		}

		chCloseCh := make(chan *amqp.Error)
		app.Queue().NotifyChannelClose(chCloseCh)

		registerGeolocationInsertHandler(ctx, app.Queue())

		<-chCloseCh
		cancel()
	}
}
