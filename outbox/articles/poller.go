package articles

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

// LikedArticlesPoller provides a mechanism for sending liked article events to
// a downstream consumer.
type LikedArticlesPoller struct {
	storage        ArticleStorage
	lastEventIndex int
	pollInterval   time.Duration

	timeout time.Duration

	conn net.Conn

	enc *json.Encoder
	dec *json.Decoder

	// Do not send new events until this specific one is acked. This also
	// means not fetching any new events until that point.
	unacked        Message
	lastEventAcked bool
}

// Message models a message that is sent over the network. It has an Id that
// should be used by the consumer for ACKs.
type Message struct {
	Id     int                 `json:"id"`
	Events []ArticleLikedEvent `json:"events"`
}

// Ack models an ACK received by the consumer.
type Ack struct {
	Id int `json:"id"`
}

func NewLikedArticlesPoller(storage ArticleStorage, pollInterval time.Duration, authority string) (*LikedArticlesPoller, error) {

	// authority is the fancy RFC 3986 name for host + port
	conn, err := net.Dial("tcp", authority)
	if err != nil {
		return nil, fmt.Errorf("could not connect to downstream: %s", err.Error())
	}

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	timeout := pollInterval - 1*time.Second

	return &LikedArticlesPoller{
		storage:        storage,
		lastEventIndex: -1,
		pollInterval:   pollInterval,
		timeout:        timeout,
		conn:           conn,
		lastEventAcked: true,
		enc:            enc,
		dec:            dec,
		unacked:        Message{},
	}, nil
}

func (p *LikedArticlesPoller) Poll() {
	ticker := time.NewTicker(p.pollInterval)
	go func() {
		cycle := 0
		for {
			cycle++
			log.Info().Int("cycle", cycle).Msg("starting new polling cycle")
			<-ticker.C

			newEvents, err := p.storage.GetArticleLikedEventsFromIndex(p.lastEventIndex)

			// staff level engineer error handling
			if err != nil {
				// fmt.Printf("error retrieving latest events, aborting: %s\n", err.Error())
				log.Error().Err(err).Msg("error retrieving latest events, aborting")
				break
			}

			// synchronously send new events. polling interval should be less
			// than downstream timeout
			err = p.sendNewEvents(newEvents)
			if err != nil {
				// fmt.Printf("did not manage to send new events: %s\n", err.Error())
				log.Error().Err(err).Msg("did not manage to send new events")
			}
		}
	}()
}

func (p *LikedArticlesPoller) sendNewEvents(events []ArticleLikedEvent) error {

	// If there are no new events, ignore this.
	if len(events) == 0 {
		log.Debug().Msg("skipping sending events, no new events")
		return nil
	}

	// fmt.Println("encoding events", events)
	log.Info().Any("events", events).Msg("encoding events")

	// Create a new message. If last event was not acked, the ACK may have
	// been lost or will come with delay. Send the same event, to ensure
	// that the ACK is generated again in case it was lost due to network
	// failure.
	//
	// Note that not using latest events plays an important role here, if
	// new events were generated. If this includes new events since an unacked
	// event was send, it gets hard to keep track of what the consumer has
	// actually persisted, and it may lead to missed events.
	var message Message
	if !p.lastEventAcked {
		message = p.unacked
	} else {
		message = Message{
			Id:     events[0].EventId,
			Events: events,
		}
	}

	err := p.enc.Encode(message)
	if err != nil {
		return fmt.Errorf("could not send events: %s", err.Error())
	}

	// fmt.Println("events encoded", events)
	log.Info().Any("message", message).Msg("message sent")

	// If the ack is not received after the timeout, we consider that the
	// consumer failed.
	p.conn.SetReadDeadline(time.Now().Add(p.timeout))

	var ack Ack
	err = p.dec.Decode(&ack)

	// Note that if we sent the same event multiple times, it doesn't matter
	// that the ACK may be delayed, nor that we may receive multiple ACKs.
	// Since from now we will be sending a new message, it will eventually
	// catch up after a few failures.
	//
	// This has the potential to lead to quite a significant delay if we
	// keep missing ACKs, if we only read from the consumer once per polling
	// interval. To mitigate that, keep reading until the timeout.
	var loopErr error
LOOP:
	for {
		// this error includes timeout
		if err != nil {
			var e *net.OpError // internal net error
			ok := errors.As(err, &e)
			if !ok {
				// this covers all errors, including decoding errors. It means
				// the event did not get ACKed properly.
				p.lastEventAcked = false
				p.unacked = message
				loopErr = err
			} else {
				if e.Timeout() {
					// This also counts as not being acked, especially the first
					// time. The difference here is that we're breaking the
					// loop.
					p.lastEventAcked = false
					p.unacked = message
					loopErr = err
					break LOOP
				}
			}
			return fmt.Errorf("failed reading ack from consumer: %s", err.Error())
		} else {

			// fmt.Println("ack received", ack)
			log.Info().Any("ack", ack).Msg("ack received")

			// Received an ACK for a previous message. We can disregard this.
			if ack.Id != message.Id {
				log.Info().Msg("received ack for different event, skipping")
				p.lastEventAcked = false
				p.unacked = message
			} else {
				// proper ACK. we should break the loop.
				p.lastEventAcked = true
				p.unacked = Message{}
				p.lastEventIndex += len(message.Events)
				break LOOP
			}
		}
	}

	if loopErr != nil {
		return fmt.Errorf("failed to ack; last received error: %s", err.Error())
	}
	return nil
}
