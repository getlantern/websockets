package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"

	"golang.org/x/net/websocket"
)

const (
	notificationsPerMinute = 23
)

var (
	connMutex sync.RWMutex
	i         = int64(0)
	conns     = make(map[int64]chan string)

	targetNotifyDuration = (60 * 1000 * 1000 / notificationsPerMinute) * time.Microsecond
)

func main() {
	flag.Parse()
	glog.Infof("Target notify duration: %v", targetNotifyDuration)

	http.Handle("/", websocket.Handler(Register))
	l, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		glog.Fatalf("Unable to listen: %v", err)
	}
	glog.Infof("Listening at %v", l.Addr())

	go notify()

	err = http.Serve(l, nil)
	if err != nil {
		glog.Fatalf("Error serving: %v", err)
	}
}

func Register(conn *websocket.Conn) {
	ch := make(chan string)
	connMutex.Lock()
	conns[i] = ch
	i += 1
	connMutex.Unlock()
	for msg := range ch {
		fmt.Fprint(conn, msg)
	}
}

func notify() {
	for {
		var ch chan string
		start := time.Now()
		connMutex.RLock()
		if i > 0 {
			ch = conns[rand.Int63n(i)]
		}
		connMutex.RUnlock()
		if ch != nil {
			glog.Info("Notifying")
			ch <- "Someone signed up with your referral code, you have a bonus!"
		} else {
			glog.Info("Noone to notify")
		}
		delta := time.Now().Sub(start)
		extraSleep := targetNotifyDuration - delta
		if extraSleep > 0 {
			glog.Infof("Sleeping %v", extraSleep)
			time.Sleep(extraSleep)
		}
	}
}
