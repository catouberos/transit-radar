package queues

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareQueues(ch *amqp.Channel) error {
	err := declareGeolocationQueues(ch)
	if err != nil {
		return err
	}

	err = declareRouteQueues(ch)
	if err != nil {
		return err
	}

	err = declareVariantQueues(ch)
	if err != nil {
		return err
	}

	return nil
}

func declareGeolocationQueues(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"geolocation", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"geolocationCreated", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                      // queue name
		"geolocation.event.created", // routing key
		"geolocation",               // exchange
		false,
		nil,
	)

	return err
}

func declareRouteQueues(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"route",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"routeUpdated", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                // queue name
		"route.event.updated", // routing key
		"route",               // exchange
		false,
		nil,
	)

	return err
}

func declareVariantQueues(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"variant", // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"variantUpdate", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                  // queue name
		"variant.event.updated", // routing key
		"variant",               // exchange
		false,
		nil,
	)

	return err
}
