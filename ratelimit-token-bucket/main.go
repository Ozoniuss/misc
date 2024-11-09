package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Simple poc function
type stringReturner func() (string, error)

// RateLimitFunction returns a function doing the same as the original function
// but with a build in rate limit.
func RateLimitFunction(f stringReturner, refillInterval time.Duration, refillSize int, bucketSize int) stringReturner {

	remaining := bucketSize
	refillTicker := time.NewTicker(refillInterval)
	bmu := &sync.Mutex{}
	go func() {
		for range refillTicker.C {
			bmu.Lock()
			remaining = min(bucketSize, remaining+refillSize)
			bmu.Unlock()
		}
	}()

	return func() (string, error) {
		bmu.Lock()
		defer bmu.Unlock()
		if remaining <= 0 {
			return "", errors.New("rate limited")
		}
		hello, err := f()
		// decrease quota
		remaining--
		return hello, err
	}
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	var printhello stringReturner = func() (string, error) {
		return "hello", nil
	}
	// Every 5 seconds, refill with 2 calls, and start with 3 calls
	f := RateLimitFunction(printhello, 5*time.Second, 2, 3)

	t := time.NewTicker(1 * time.Second)
	for range t.C {
		n := time.Now()
		val, err := f()
		if err != nil {
			fmt.Printf("[%s] %s\n", n.Format(time.TimeOnly), err.Error())
		} else {
			fmt.Printf("[%s] %s\n", n.Format(time.TimeOnly), val)
		}
	}
}

// func main() {

// 	t := time.NewTicker(1 * time.Second)
// 	for range t.C {
// 		fmt.Println("ce ma?")
// 	}

// 	select {
// 	case v, ok := <-t.C:
// 		fmt.Println("lol", v, ok)
// 	default:
// 		fmt.Println("hai ca esti parlit")
// 	}

// 	time.Sleep(2 * time.Second)

// 	select {
// 	case v, ok := <-t.C:
// 		fmt.Println("lol", v, ok)
// 	default:
// 		fmt.Println("hai ca esti parlit")
// 	}

// }
