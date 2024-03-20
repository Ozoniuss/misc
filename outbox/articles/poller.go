package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

// LikedArticlesPoller provides a mechanism for sending liked article events to
// a downstream consumer.
type LikedArticlesPoller struct {
	storage        ArticleStorage
	lastEventIndex int
	pollInterval   time.Duration

	conn net.Conn
}

func NewLikedArticlesPoller(storage ArticleStorage, pollInterval time.Duration, authority string) (*LikedArticlesPoller, error) {

	// authority is the fancy RFC 3986 name for host + port
	conn, err := net.Dial("tcp", authority)
	if err != nil {
		return nil, fmt.Errorf("could not connect to downstream: %s", err.Error())
	}

	return &LikedArticlesPoller{
		storage:        storage,
		lastEventIndex: -1,
		pollInterval:   pollInterval,
		conn:           conn,
	}, nil
}

func (p *LikedArticlesPoller) poll() {
	ticker := time.NewTicker(p.pollInterval)
	go func() {
		for {
			<-ticker.C
			newEvents, err := p.storage.GetArticleLikedEventsFromIndex(p.lastEventIndex)

			// staff level engineer error handling
			if err != nil {
				fmt.Printf("error retrieving latest events, aborting: %s\n", err.Error())
				break
			}

			// synchronously send new events. polling interval should be less
			// than downstream timeout
			err = p.sendNewEvents(newEvents)
			if err != nil {
				fmt.Printf("did not manage to send new events: %s\n", err.Error())
			}
		}
	}()
}

func (p *LikedArticlesPoller) sendNewEvents(events []ArticleLikedEvent) error {
	enc := gob.NewEncoder(p.conn)
	err := enc.Encode(events)
	if err != nil {
		return fmt.Errorf("could not send events: %s", err.Error())
	}

	var ack bool
	dec := gob.NewDecoder(p.conn)
	err = dec.Decode(&ack)
	if err != nil {
		return fmt.Errorf("failed reading ack from consumer: %s", err.Error())
	}

	// This is not particularly necessary, just an extra check for making sure
	// that the consumeD send the intended ack.
	if ack {
		p.lastEventIndex += len(events)
	}
	return nil
}
