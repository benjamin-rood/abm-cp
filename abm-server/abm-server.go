package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/benjamin-rood/goio"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	om       = make(chan goio.OutMsg)
	phase    = make(chan struct{})
	view     = make(chan render.Viewport)
	ctxt     = make(chan abm.Context)
	mdl      = make(chan struct{}) // to send close to running processes on socket failure.
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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

func reader(conn *websocket.Conn, ws chan struct{}, model chan struct{}, ctxt chan<- abm.Context, view chan<- render.Viewport) {
	defer conn.Close()
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var msgIn interface{}
	var rawIn goio.InMsg
	var err error
	for {
		select {
		case <-ws:
			close(model)
			return
		default:
			if err = conn.ReadJSON(&rawIn); err != nil {
				networkError(err, ws)
			}
			data := []byte(rawIn.Data)
			switch {
			case rawIn.Type == "context":
				msgIn = abm.Context{}
				if err = json.Unmarshal(data, &msgIn); err != nil {
					dataError(err, ws)
				}
				ctxt <- msgIn.(abm.Context)
			case rawIn.Type == "viewport":
				msgIn = render.Viewport{}
				if err = json.Unmarshal(data, &msgIn); err != nil {
					dataError(err, ws)
				}
				view <- msgIn.(render.Viewport)
			}
		}
	}
}

func writer(conn *websocket.Conn, ws chan struct{}, om <-chan goio.OutMsg) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		conn.Close()
	}()
	for {
		select {
		case msg := <-om:
			fmt.Println(len(msg.Data.(render.DrawList).CPP))
			fmt.Println(len(msg.Data.(render.DrawList).VP))
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("writer: failed to WriteJSON:")
				networkError(err, ws)
			}
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writer: ping msg fail:", err)
				networkError(err, ws)
			}
		case <-ws:
			// clean up here, then...
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			networkError(err, nil)
		}
		return
	}
	ws := make(chan struct{}) //	websocket signalling channel
	abm.InitModel(abm.DemoContext, abm.DemoEnvironment, om, view, phase)
	go writer(conn, ws, om)
	reader(conn, ws, mdl, ctxt, view)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/ws", serveWs)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
