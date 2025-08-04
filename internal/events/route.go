package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/dto"
)

func registerRouteUpsertHandler(ctx context.Context, app *base.App) error {
	deliveries, err := app.Queue().Consume(
		"routeUpdated", // queue
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
				logger.With("queue", "routeUpdated")

				params := &dto.RouteUpsert{}
				err := json.Unmarshal(delivery.Body, params)
				if err != nil {
					logger.Error("Error unmarshal queue data", "error", err)
				}

				_, err = app.CreateOrUpdateRoute(ctx, params)
				if err != nil {
					logger.Error("Error creating route", "error", err)
				}

				if err := delivery.Ack(false); err != nil {
					logger.Error("Error acknowledging message", "error", err)
				}
			}
		}
	}()

	return nil
}
