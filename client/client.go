package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

const (
	numDialers = 100
	origin     = "http://localhost/"
	url        = "ws://localhost:12345/"
)

var (
	finished sync.WaitGroup
)

func main() {
	numConnections, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to parse num connections: %v", err)
	}
	fmt.Printf("Launching %v connections\n", numConnections)
	dialed := int64(0)
	var doneDialing sync.WaitGroup
	doneDialing.Add(numDialers)
	for j := 0; j < numDialers; j++ {
		go func() {
			for i := 0; i < numConnections/numDialers; i++ {
				ws, err := websocket.Dial(url, "", origin)
				if err != nil {
					log.Printf("Unable to dial: %v", err)
					continue
				}
				finished.Add(1)
				atomic.AddInt64(&dialed, 1)
				go read(ws)
			}
			doneDialing.Done()
		}()
	}
	doneDialing.Wait()
	fmt.Printf("Launched %v connections\n", atomic.LoadInt64(&dialed))
	finished.Wait()
}

func read(ws *websocket.Conn) {
	defer finished.Done()
	var msg = make([]byte, 512)
	for {
		n, err := ws.Read(msg)
		if err != nil {
			return
		}
		fmt.Println(string(msg[:n]))
	}
}
