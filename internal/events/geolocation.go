package events

import (
	"context"
	"log"

	"github.com/catouberos/geoloc/base"
)

func registerGeolocationInsertHandler(app *base.App) error {
	ctx := context.Background()

	msgs, err := app.AMQP.ConsumeWithContext(
		ctx,
		"geolocation.insert", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
			}
		}()

		// run until canclled
		<-ctx.Done()
	}(ctx)

	return nil
}
