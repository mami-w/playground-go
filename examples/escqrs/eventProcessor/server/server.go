package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mami-w/playground-go/examples/escqrs/dronescommon"
	repo2 "github.com/mami-w/playground-go/examples/escqrs/eventProcessor/repo"
	"github.com/streadway/amqp"
	"log"
)

func NewServer(url string) (router *mux.Router){

	router = mux.NewRouter()

	initRoutes(router)

	telemetryChannel := make(chan dronescommon.TelemetryUpdatedEvent)
	positionChannel := make(chan dronescommon.PositionChangedEvent)

	repo := repo2.InitRepository()

	dequeEvents(telemetryChannel, positionChannel, url)
	consumeEvents(telemetryChannel, positionChannel, repo)

	return router
}

func initRoutes(router *mux.Router) {
	// anything to do here? probably not, just keep running...
}

func dequeEvents(telemetryChannel chan dronescommon.TelemetryUpdatedEvent,
	positionChannel chan dronescommon.PositionChangedEvent, url string) {

	conn, err := amqp.Dial(url)
	fatalError(err)

	ch, err := conn.Channel()
	fatalError(err)

	positionsQ, _ := ch.QueueDeclare("positions", false, false, false, false, nil)
	telemetryQ, _ := ch.QueueDeclare("telemetry", false, false, false, false, nil)

	positionsIn, _ := ch.Consume(positionsQ.Name, "", true, false, false, false, nil)
	telemetryIn, _ := ch.Consume(telemetryQ.Name, "", true, false, false, false, nil)

	go func() {
		for {
			select {
				case positionRaw := <-positionsIn:
					dispatchPosition(positionRaw, positionChannel)
				case telemetryRaw := <-telemetryIn:
					dispatchTelemetry(telemetryRaw, telemetryChannel)
			}
		}
	}()
}

func dispatchPosition(alertRaw amqp.Delivery, out chan dronescommon.PositionChangedEvent) {
	var event dronescommon.PositionChangedEvent
	err := json.Unmarshal(alertRaw.Body, &event)
	if err != nil {
		out <- event
	} else {
		log.Print("could not dequeue position event")
	}
}

func dispatchTelemetry(alertRaw amqp.Delivery, out chan dronescommon.TelemetryUpdatedEvent) {
	var event dronescommon.TelemetryUpdatedEvent
	err := json.Unmarshal(alertRaw.Body, &event)
	if err == nil {
		out <- event
	} else {
		log.Print("could not dequeue telemetry event")
	}
}

func consumeEvents(telemetryChannel chan dronescommon.TelemetryUpdatedEvent,
	positionChannel chan dronescommon.PositionChangedEvent, repo *repo2.Repository) {

		go func() {
			for {
				select {
					case position := <- positionChannel:
						repo.Save(position)
					case telemetry := <- telemetryChannel:
						repo.Save(telemetry)
				}
			}
		}()
}

func fatalError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

