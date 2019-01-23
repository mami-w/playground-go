package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mami-w/playground-go/examples/escqrs/commandHandler/commandHandlers"
	"github.com/mami-w/playground-go/examples/escqrs/commandHandler/types"
	"github.com/streadway/amqp"
	"log"
) // https://github.com/gorilla/mux

// NewServer configures and returns a Server.
func NewServer(url string)  *mux.Router {

	mx := mux.NewRouter()

	positionDispatcher := buildDispatcher("positions", url)
	telemetryDispatcher := buildDispatcher("telemetry", url)

	initRoutes(mx, telemetryDispatcher, positionDispatcher)

	return mx
}

func initRoutes(mx *mux.Router, telemetryDispatcher types.QueueDispatcher, positionDispatcher types.QueueDispatcher) {
	mx.HandleFunc("/api/cmds/telemetry", commandHandlers.AddTelemetryHandler(telemetryDispatcher)).Methods("POST")
	//mx.HandleFunc("/api/cmds/positions", addPositionHandler(positionDispatcher)).Methods("POST")
}

func buildDispatcher(dispatcherType string, url string) (dispatcher types.QueueDispatcher) {

	fmt.Printf("\nUsing URL (%s) for Rabbit.\n", url)

	conn, err := amqp.Dial(url)
	handleError(err)

	ch, err := conn.Channel()
	handleError(err)

	q, err := ch.QueueDeclare(
		dispatcherType, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	handleError(err)

	dispatcher = newDispatcher(ch, q.Name, false)
	return dispatcher
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}