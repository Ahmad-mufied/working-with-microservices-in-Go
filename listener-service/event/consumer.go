package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

/*
The NewConsumer function is called with an AMQP connection object as input. It creates a new Consumer struct and calls the setup method
to set up the connection and declare the exchange.
*/

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil

}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

/*
The Listen method is called with a slice of topics as input. It creates a new channel, declares a random queue,
and binds the queue to the specified topics.The method then starts consuming messages from the queue and for each message,
it unmarshals the payload JSON into a Payload struct and calls the handlePayload function in a separate goroutine.
*/

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

/*
	The handlePayload function checks the payload name and performs different actions based on the name.
	If the name is "log" or "event", it calls the logEvent function to log the payload. If the name is "auth",
	it performs authentication (not implemented in the code snippet). Otherwise, it also calls the logEvent function.
*/

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
	// authenticate

	// you can have as many cases as you want, as long as you write the logic
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

/*
	The logEvent function marshals the payload into JSON, creates an HTTP POST request to a logging service,
	sends the request, and checks the response status code.
*/

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
