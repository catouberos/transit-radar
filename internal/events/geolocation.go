package events

import (
	"context"
	"log"
	"time"

	"github.com/catouberos/geoloc/base"
)

type GeolocationInsert struct {
	Degree    float32   `json:"deg"`
	Latitude  float32   `json:"lat"`
	Longitude float32   `json:"lng"`
	Speed     float32   `json:"spd"`
	VehicleID int64     `json:"vID"`
	RouteID   int64     `json:"rID"`
	Timestamp time.Time `json:"dts"`
}

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
