package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// A simple outbox pattern example. Essentially:
//
// - Have some object on which you can do atomic operations on ("aggregate").
// For simplicity this will be just an int, which is the current "state" of the
// main program.
// - When doing an operation save both the new state and an event of what was
// changed. The program allows adding a value to the number, so the addend is
// the event.
// - The events must be saved locally (in memory), and also be sent to the
// listener.
// A new thread will poll for the queue events every few miliseconds and send
// the new events through a channel.
// - When the listener acks the event(s), remove the acked event(s) from the
// queue.
//
//  Constraints:
// - All events have an ID. It is the responsibility of the main program to
// keep track of that ID and generate unique IDs.
// - Receiving an event and storing it in the event array (slice) is an atomic
// operation. It will not be possible for the polling thread to interfere just
// between receiving an event and storing it in the queue, nor just modifying
// the state without storing the event.
// - Events are received by the main program in the order they were sent.
// - For the sake of this demo it is assumed that the program and thread do
// not crash on their own. The relevant failures to the exercise are modeled
// through code.
//
// Failure scenarios:
//
// - Event delivery is guaranteed to happen at least once. The listener may
// fail to ack (which is modeled by the program), in which case the event will
// be re-sent again. It is the responsibility of the listener to keep track of
// which events were received and apply them accordingly.

func listener(events <-chan [2]int) {
	internalState := 0
	lastReceivedEvent := -1
	for val := range events {
		// the sender didn't receive the ack
		if lastReceivedEvent == val[0] {
			fmt.Println("receivied an existing event")
			continue
		}
		// Note that in theory, this would be implemented as an atomic operation
		internalState += val[1]
		lastReceivedEvent = val[0]
		fmt.Printf("got event %d, current state %d\n", val[1], internalState)
	}
	_ = lastReceivedEvent // hacks to avoid compile time checks
	fmt.Println("channel was closed")
}

// Database mocks a real database in the application's memory. The state is
// meant to represent a set of columns representing the aggregate's state,
// potentially with a set of constraints. Events is meant to represent a
// separate table (called the outbox table) storing all the unpublished
// domain events (or all events, and marking the unpublished ones).
//
// Note that it is not necessarily the database's job to ensure the aggregate
// constraints are satisfied after each operation. That falls on the aggregate's
// internal logic.
type Database struct {
	// Implementation detail that protects the database from concurrent access
	// between the aggregate operation and the poller. This is specific to the
	// in-memory implementation.
	mtx *sync.Mutex

	state  int
	events []DomainEvent

	// This can be computed easily from the database. It's useful for the
	// consumer to distinguish between events that were already processed.
	lastEventId int
}

type AggregateType struct {
	dbConn *Database
}
type DomainEvent struct {
	id   int
	diff int
}

// updateValue is an operation that can be performed on the aggregate.
func (a *AggregateType) updateValue(new int) {
	// Simple aggregate internal logic
	diff := new - int(a.dbConn.state)

	// Other bussiness logic, e.g. for determining that the constraints are
	// satified
	// ...

	// Atomic database transaction, which writes the new state.
	a.dbConn.mtx.Lock()
	defer a.dbConn.mtx.Unlock()
	a.dbConn.state = new
	a.dbConn.events = append(a.dbConn.events, DomainEvent{
		id:   a.dbConn.lastEventId + 1,
		diff: diff,
	})
	a.dbConn.lastEventId++
}

// pollDatabaseForUnpublishedEvents represents the message relay which tries to
// send the unpublished domain events to the consumers. This is done via a
// poll-based strategy, rather than a push-based one.
//
// A failure scenario is simulated within the message relay using randomization.
func pollDatabaseForUnpublishedEvents(dbConn *Database, events chan<- []DomainEvent) {
	for {
		time.Sleep(1 * time.Second)
		sendUnpublishedEvents(dbConn, events)
	}
}

// sendUnpublishedEvents attempts to retrieve all the unpublished events from
// the database and send them to the downstream consumer.
func sendUnpublishedEvents(dbConn *Database, events chan<- []DomainEvent) {
	dbConn.mtx.Lock()
	defer dbConn.mtx.Unlock()

	// No new events to publish.
	if len(dbConn.events) == 0 {
		return
	}

	// Simulate periodic failure, that is either the event bus failing after
	// publishing but before marking the events as being published (removing
	// them from the database) or not receiving an ACK after a certain period
	// of time.
	r := rand.Intn(100)
	if r > 95 {
		// Events were sent but not removed, which is the unsuccessful
		// scenario.
		events <- dbConn.events
	} else {
		// Events were both sent and removed, which is the successful
		// scenario.
		events <- dbConn.events
		dbConn.events = []DomainEvent{}
	}
}

func main() {
	events := make(chan [2]int, 100)
	go listener(events)

	// Mocks an internal state of some aggregate. Think of an aggregate as an
	// abstraction holding a set of constraints that have to be satisfied. All
	// operations performed on the aggregate are atomic (they either success or
	// fail) and ensure its contraints are satisified.
	//
	// In this instance, the aggregate is simply mocked by an integer, and the
	// only constraint is that the state is always represented by an integer.
	state := AggregateType(0)

	// Events that were emitted by the aggregate when its state changed.
	// When the aggregate is changed, it emits an event with the difference
	// between the new value and the old value.
	eventsArray := make([]int, 0)

	eventId := 0
	mtx := &sync.Mutex{}

	go sendUpdates(&eventsArray, &eventId, mtx, events)

	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("write a number or else: ")
		s.Scan()
		num, err := strconv.Atoi(s.Text())
		if err != nil {
			fmt.Println("u fked up")
		}

		incrementStateAtomically(num, &state, &eventsArray, &eventId, mtx)
	}

}

// incrementstateatomically either increments the state and stores the event or
// fails altogeher
func incrementStateAtomically(input int, state *int, eventsArray *[]int, eventId *int, mtx *sync.Mutex) {
	r := rand.Intn(100)

	// throw a random ass message
	if r > 50 {
		fmt.Println("the system fked up")
		return
	}

	// Avoid race conditions with the "event pusher".
	mtx.Lock()
	defer mtx.Unlock()
	*state += input
	(*eventsArray) = append((*eventsArray), input)
	*eventId++
}
