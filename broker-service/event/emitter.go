package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

// setup sets up a channel for communication with RabbitMQ.
// It returns an error if there is an issue creating the channel or declaring the exchange.
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchange(channel)
}

// Push publishes an event to the RabbitMQ exchange.
// It takes the event message and severity level as inputs.
// It returns an error if there is an issue creating the channel or publishing the event.
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewEventEmitter creates a new instance of the Emitter struct and calls the setup method.
// It takes a connection to RabbitMQ as input.
// It returns the new Emitter instance and an error if there is an issue setting up the channel.
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
