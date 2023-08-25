package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		_ = c.Close()
	}()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","value":"outcomes"}`)); err != nil {
		log.Println("write:", err)
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","value":"leaderboard"}`)); err != nil {
		log.Println("write:", err)
		return
	}

	select {
	case <-interrupt:
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write last:", err)
			return
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
		return
	case <-done:
		return
	}
}
