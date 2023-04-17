package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

var DEPTH int = 1

// pow2 raises two to the power of n, returning an int result.
func pow2(n int) int {
	ret := 1
	for i := 1; i <= n; i++ {
		ret *= 2
	}
	return ret
}

// log2 returns the depth of a given goroutine.
func log2(n int) int {
	depth := 0
	for n > 0 {
		n = n / 2
		depth++
	}
	return depth
}

func runInterruptible(idx int, ctx context.Context, kill chan struct{}) {

	// guard against the main thread.
	if idx == 0 {
		return
	}
	if log2(idx) < DEPTH {
		childctx, cancel := context.WithCancel(context.Background())
		forkWithCtx(idx, childctx)

		// This will kill the children once the goroutine is killed.
		defer cancel()
	}
	fmt.Printf("[%d] running\n", idx)

	// This will notify the main goroutine that the children had died.
	defer func() {
		done <- idx
	}()

	// Listen for signals.
	select {
	case <-kill:
		fmt.Printf("[%d] killed\n", idx)
		return
	case <-ctx.Done():
		fmt.Printf("[%d] cancelled\n", idx)
		return
	}
}

// forkWithCtx creates two children with a given context and runs them.
func forkWithCtx(idx int, ctx context.Context) {
	go runInterruptible(2*idx, ctx, kill[2*idx])
	go runInterruptible(2*idx+1, ctx, kill[2*idx+1])
}

func readInput() {
	var in int
	for {
		_, err := fmt.Scan(&in)
		if err != nil {
			fmt.Printf("Invalid input: %s\n", err.Error())
			continue
		}
		if in < 0 || in > pow2(DEPTH)-1 {
			fmt.Printf("Id must be between 1 and %d\n", pow2(DEPTH)-1)
			continue
		}
		if in == 0 {
			fmt.Println("Cannot stop main directly. Stop all the child goroutines instead.")
			fmt.Printf("Goroutines still alive: %v\n", getAlive())
			continue
		}
		_, alive := alive[in]
		if alive {
			kill[in] <- struct{}{}
		} else {
			fmt.Printf("[%d] already dead\n", in)
		}
	}
}

// Signals the main process that all processes had been terminated.
var done chan int

// Channels that signal each individual process that it was shut down.
// kill[i] is assigned to process with number i, starting from 1.
var kill []chan struct{}

// closed determines which goroutines had already been killed.
var alive map[int]struct{}

// Used to write which processes are still alive.
var getalive strings.Builder

func getAlive() string {
	// Remove old buffer.
	getalive.Reset()
	start := true
	for k := range alive {
		// 	Ensures to not add a comma at the end.
		if !start {
			getalive.WriteByte(',')
		}
		getalive.WriteString(strconv.FormatInt(int64(k), 10))
		start = false
	}
	return getalive.String()
}

func main() {

	fmt.Println("What is the depth of the three?")
	var depth int = 0
	for {
		_, err := fmt.Scan(&depth)
		if err != nil {
			fmt.Printf("Invalid depth: %s", err.Error())
		} else if depth < 1 {
			fmt.Println("Depth must be positive.")
		} else {
			break
		}
	}
	DEPTH = depth

	done = make(chan int)
	alive = make(map[int]struct{}, pow2(depth))
	getalive = strings.Builder{}

	kill = make([]chan struct{}, pow2(depth))
	for i := 1; i <= pow2(DEPTH)-1; i++ {
		kill[i] = make(chan struct{})
		alive[i] = struct{}{}
	}

	fmt.Println("[0] main running")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	forkWithCtx(0, ctx)

	go readInput()

	// Waits for all goroutines to die. This is an unbuffered channel, so a
	// goroutine stopping happens simultaneously with main recognizing it.
	// This will update the map with the goroutine that died.
	for i := 0; i < pow2(DEPTH)-1; i++ {
		idx := <-done
		delete(alive, idx)
	}

	fmt.Println("main stopped")
}
