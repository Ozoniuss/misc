package main

import (
	"fmt"
	"sync"
	"time"

	log "github.com/Ozoniuss/stdlog"
)

var took = make(chan string, 1)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	elapsedstr := fmt.Sprintf("%v", elapsed)
	took <- elapsedstr
}

func writealot(stdwriter func(...any)) {
	defer timeTrack(time.Now())
	w := &sync.WaitGroup{}
	w.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			for j := 0; j < 1000; j++ {
				stdwriter("something")
			}
			w.Done()
		}(i)
	}
	w.Wait()
}

func logalot() {
	writealot(func(a ...any) {
		log.Infoln(a)
	})
}

func printalot() {
	writealot(func(a ...any) {
		// Just to have the same format
		fmt.Println("[info] 2023/02/17 20:03:44", a)
	})
}

func main() {
	printalot()
	printalotTook := <-took
	logalot()
	logalotTook := <-took

	fmt.Printf("printalot: %s, logalot: %s", printalotTook, logalotTook)
}
