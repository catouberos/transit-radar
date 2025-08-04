package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/base"
)

func registerVariantUpsertHandler(ctx context.Context, app *base.App) error {
	deliveries, err := app.Queue().Consume(
		"variantUpdated", // queue
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
				logger.With("queue", "variantUpdated")

				params := &dto.VariantByRouteEbmsIDUpsert{}
				err := json.Unmarshal(delivery.Body, params)
				if err != nil {
					logger.Error("Error unmarshal queue data", "error", err)
				}

				_, err = app.CreateOrUpdateVariantByRouteEbmsID(ctx, params)
				if err != nil {
					logger.Error("Error creating variant", "error", err)
				}

				if err := delivery.Ack(false); err != nil {
					logger.Error("Error acknowledging message", "error", err)
				}
			}
		}
	}()

	return nil
}
