package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/base"
)

func registerGeolocationInsertHandler(ctx context.Context, app *base.App) error {
	deliveries, err := app.Queue().Consume(
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
				logger := slog.New(slog.Default().Handler())
				logger.With("queue", "geolocationCreated")

				params := &dto.GeolocationByRouteIDAndPlateAndBoundInsert{}
				err := json.Unmarshal(delivery.Body, params)
				if err != nil {
					logger.Error("Error unmarshal queue data", "error", err)
				}

				_, err = app.CreateGeolocationByRouteIDAndPlateAndBound(ctx, params)
				if err != nil {
					logger.Error("Error creating geolocation", "error", err)
				}

				if err := delivery.Ack(false); err != nil {
					logger.Error("Error acknowledging message", "error", err)
				}
			}
		}
	}()

	return nil
}
