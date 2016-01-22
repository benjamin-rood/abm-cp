package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/goio"
	"golang.org/x/net/websocket"
)

func networkError(err error, c chan struct{}) {
	log.Println(err)
	if err.Error() == "use of closed network connection" {
		close(c)
	}
}

func modelError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

func dataError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

type abmHub map[string]abm.Client

var hub abmHub

func serverSession(ws *websocket.Conn) {
	uid := getUIDString()
	c := abm.Client{
		UID:    uid,
		Model:  abm.NewModel(),
		Active: true,
		Stamp:  time.Now(),
	}
	hub[uid] = c
	go wsReader(ws, c.Im)
	go wsWriter(ws, c.Om)
	watcher(&c)
}

func watcher(c *abm.Client) {
	defer func() {
		c.Active = false
		c.Stamp = time.Now()
	}()
	select {
	case <-c.Q:
		// send final statistics
		// clean up
		return
	}
}

func wsReader(ws *websocket.Conn, in chan<- goio.InMsg) {
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		var msg goio.InMsg
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Println("wsReader: Can't receive!")
		}
		in <- msg
	}
}

func wsWriter(ws *websocket.Conn, out <-chan goio.OutMsg) {
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-out:
			websocket.JSON.Send(ws, msg)
		}
	}
}

func getUIDString() string {
again:
	id := randomdata.SillyName()
	_, exists := hub[id]
	if exists {
		goto again
	}
	return id
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/ws", websocket.Handler(serverSession))
	http.ListenAndServe(":8080", nil)
}
