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

func sum(nums []int) int {
	s := 0
	for i := 0; i < len(nums); i++ {
		s += nums[i]
	}
	return s
}

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

// simulates either receiving or not receiving an ack from the listener, as in,
// failing is the same as not receiving an ack
func publishOrFail(eventsArray *[]int, eventId *int, mtx *sync.Mutex, events chan<- [2]int) {

	// no events to publish
	if len(*eventsArray) == 0 {
		return
	}

	// For the sake of this example this cannot happen while the state is being
	// incremented so there are no race conditions.
	mtx.Lock()
	defer mtx.Unlock()

	// we assume all events were sent at once. this is functionally equivalent
	// to sending the sum
	eventCompensated := sum(*eventsArray)

	r := rand.Intn(100)
	// random ass error
	if r > 50 {
		// In this case we send the event but didn't receive an ack. do not
		// empty the q
		events <- [2]int{*eventId, eventCompensated}
	} else {
		// empty the events array
		events <- [2]int{*eventId, eventCompensated}
		*eventsArray = []int{}
	}

}

func sendUpdates(eventsArray *[]int, eventId *int, mtx *sync.Mutex, events chan<- [2]int) {
	for {
		time.Sleep(500 * time.Millisecond)
		publishOrFail(eventsArray, eventId, mtx, events)
	}
}

func main() {
	events := make(chan [2]int, 100)
	go listener(events)

	state := 0

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
