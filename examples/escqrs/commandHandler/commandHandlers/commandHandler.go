package commandHandlers

import (
	"encoding/json"
	"github.com/mami-w/playground-go/examples/escqrs/commandHandler/types"
	"github.com/mami-w/playground-go/examples/escqrs/dronescommon"
	"net/http"
	"io/ioutil"
	"time"
)

func AddTelemetryHandler(dispatcher types.QueueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newTelemetryCommand types.TelemetryCommand
		err := json.Unmarshal(payload, &newTelemetryCommand)
		if err != nil {
			http.Error(w, "failed", http.StatusBadRequest)
			return
		}

		if !newTelemetryCommand.IsValid() {
			http.Error(w, "todo", http.StatusBadRequest)
			return
		}

		evt := dronescommon.TelemetryUpdatedEvent{
			DroneID:newTelemetryCommand.DroneID,
			RemainingBattery:newTelemetryCommand.RemainingBattery,
			Uptime:newTelemetryCommand.Uptime,
			CoreTemp:newTelemetryCommand.CoreTemp,
			ReceivedOn:time.Now().UnixNano(),
		}
		dispatcher.DispatchMessage(evt)

		w.WriteHeader(http.StatusOK)
		// todo: what goes into body?
	}
}
