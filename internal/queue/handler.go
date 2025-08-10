package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/base"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func geolocationInsertHandler(app *base.App) rabbitmq.Handler {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		params := &dto.GeolocationByRouteIDAndPlateAndBoundInsert{}
		err := json.Unmarshal(delivery.Body, params)
		if err != nil {
			slog.Error("Error unmarshal queue data", "error", err)
			return rabbitmq.NackDiscard
		}

		_, err = app.CreateGeolocationByRouteIDAndPlateAndBound(ctx, params)
		if err != nil {
			slog.Error("Error creating geolocation", "error", err)
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}

func routeUpsertHandler(app *base.App) rabbitmq.Handler {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		params := &dto.RouteUpsert{}
		err := json.Unmarshal(delivery.Body, params)
		if err != nil {
			slog.Error("Error unmarshal queue data", "error", err)
			return rabbitmq.NackDiscard
		}

		_, err = app.CreateOrUpdateRoute(ctx, params)
		if err != nil {
			slog.Error("Error creating route", "error", err)
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}

func variantUpsertHandler(app *base.App) rabbitmq.Handler {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		params := &dto.VariantByRouteEbmsIDUpsert{}
		err := json.Unmarshal(delivery.Body, params)
		if err != nil {
			slog.Error("Error unmarshal queue data", "error", err)
			return rabbitmq.NackDiscard
		}

		_, err = app.CreateOrUpdateVariantByRouteEbmsID(ctx, params)
		if err != nil {
			slog.Error("Error creating variant", "error", err)
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}

func stopImportHandler(app *base.App) rabbitmq.Handler {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		params := &dto.StopImport{}
		err := json.Unmarshal(delivery.Body, params)
		if err != nil {
			slog.Error("Error unmarshal stop data", "error", err)
			return rabbitmq.NackDiscard
		}

		err = app.ImportStop(ctx, params)
		if err != nil {
			slog.Error("Error importing stop", "error", err)
			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	}
}
