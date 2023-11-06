package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares an exchange named "logs_topic" with the type "topic".
// The exchange is set to be durable, not auto-deleted, not internal, and with no-wait option.
// No additional arguments are provided.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      //type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguments?
	)
}

// declareRandomQueue declares a random queue with an empty name.
// The queue is set to be non-durable, not deleted when unused, exclusive, not no-wait, and no additional arguments are provided.
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name?
		false, //durable?
		false, //delete when unused?
		true,  //exclusive?
		false, //no-wait?
		nil,   // arguments?
	)
}
