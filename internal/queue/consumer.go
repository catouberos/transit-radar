package queue

import (
	"log/slog"
	"sync"

	"github.com/catouberos/transit-radar/internal/base"
	"github.com/wagslane/go-rabbitmq"
)

type consumerHandler struct {
	conn     *rabbitmq.Conn
	mu       sync.Mutex
	handlers []*handler
}

type handler struct {
	consumer *rabbitmq.Consumer
	handler  rabbitmq.Handler
}

func NewConsumerHandler(conn *rabbitmq.Conn, app *base.App) *consumerHandler {
	handler := &consumerHandler{
		conn:     conn,
		mu:       sync.Mutex{},
		handlers: []*handler{},
	}

	handler.AddConsumer("geolocationCreated", "geolocation.event.created", "geolocation", geolocationInsertHandler(app))
	handler.AddConsumer("routeUpdated", "route.event.updated", "route", routeUpsertHandler(app))
	handler.AddConsumer("variantUpdated", "variant.event.updated", "variant", variantUpsertHandler(app))

	return handler
}

func (h *consumerHandler) AddConsumer(queueName string, routingKey string, exchangeName string, handlerFn rabbitmq.Handler) error {
	consumer, err := rabbitmq.NewConsumer(
		h.conn,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = append(h.handlers, &handler{
		consumer: consumer,
		handler:  handlerFn,
	})

	return nil
}

func (h *consumerHandler) Start() {
	wg := &sync.WaitGroup{}

	for _, handler := range h.handlers {
		wg.Add(1)

		go func() {
			err := handler.consumer.Run(handler.handler)
			if err != nil {
				slog.Error("Consumer encountered an error", "error", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func (h *consumerHandler) Stop() {
	for _, handler := range h.handlers {
		handler.consumer.Close()
	}
}
