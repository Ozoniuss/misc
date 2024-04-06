package main

import (
	"fmt"
	"net/http"
	"time"
)

// 18889
func main() {

	// model event frequency in time
	eventsSeconds := []int{4, 6, 10, 8, 12, 13, 6, 3, 9, 3}
	// make a request every millisecond
	t := time.NewTicker(1 * time.Millisecond)
	tsec := time.NewTicker(1 * time.Second)

	var req, err = http.NewRequest("POST", "http://localhost:18889/articles/1/like", nil)
	if err != nil {
		panic(err)
	}

	c := 1
	iteration := 0
	shouldProcess := true
	for {
		if iteration >= len(eventsSeconds) {
			break
		}
		select {
		case <-t.C:
			if c > eventsSeconds[iteration] {
				shouldProcess = false
			}
			if shouldProcess {
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(res.Status)
				c++
			}
		case <-tsec.C:
			shouldProcess = true
			iteration++
			c = 1
		}
	}
}
