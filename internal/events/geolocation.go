package events

import (
	"context"
	"log"
	"log/slog"

	"github.com/catouberos/geoloc/internal/queues"
)

func registerGeolocationInsertHandler(ctx context.Context, queue *queues.Client) error {
	deliveries, err := queue.Consume(
		"geolocationCreated", // queue
	)
	if err != nil {
		slog.Error("Error creating consumer", "error", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-deliveries:
				log.Printf("Received a message: %s", delivery.Body)
				if err := delivery.Ack(false); err != nil {
					slog.Error("error acknowledging message: %s\n", "error", err)
				}
			}
		}
	}()

	return nil
}
