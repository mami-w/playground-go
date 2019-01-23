package main

import (
	"github.com/mami-w/playground-go/examples/escqrs/commandHandler/service"
	"net/http"
)

func main() {

	url := "amqp://guest:guest@0.0.0.0:5672"
	svr := service.NewServer(url)
	//http.Handle("/", svr)

	http.ListenAndServe(":8000", svr)
}

/*************************

	1. Event sourcing:

		f(state_1, event, ....) = state_2

		* Idempotent
			(deterministic)
		* Isolated
			ex: don't use stock data/ clock, make it an event
		* Replayable & Recoverable
		* Potential to generate lots of data

	2. Eventual Consistency

		* favors reliability & scale over instantaneous consistency

	3. CQRS

Command query responsibility segregation


Stimulus
	|
	\/
Command handler (code)
	|
	\/
Event
	|
	\/
Message queue/ event store (perm) (publish-subscribe)(Kafka) (typically something like a message queue. Needs to be perm for replay)
	|
	\/
Event processor (code)
	|
	\/
Data store(view store) (there could be a caching layer here)
	/\
	|
Query


Udi Dahan:

1. query are read only
2. view store is just a cache. Needs some reconciliation mechanism with ground truth. Does not need to be relational DB
3. Commands are a different unit of work than let's say a "screen". "Screen" is update users data, command is update married status (i.e. finer granule)
4. Commands don't need to be processed immediately, they can be queued.
5. Commands can fail on server side
6. Command handlers are independent. IMplement them as independent processes (ideally)
7. Event sourcing is only one implementation method of command hanling (domain model, transaction script, whatever)
8. Each command handler needs only simple entity to process and persist request.
9. I.e. DB could be key-value store
10. Aggregate entities hold id's of children (how does the aggregate get triggered?)
11. Only after the command has been accepted and processed and permanetn storage updated, THEN an event is published.
12. Publishing of event needs to be in same transaction as commit to permanent store.

		
	***********************************/