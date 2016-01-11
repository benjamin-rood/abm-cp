package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
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
	writeWait = 50 * time.Millisecond

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
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func reader(ws *websocket.Conn, ctxt chan<- abm.Context, view chan<- render.Viewport) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var msgIn interface{}
	var rawIn goio.InMsg
	var err error
	for {
		if err = ws.ReadJSON(&rawIn); err != nil {
			break
		}
		data := []byte(rawIn.Data)
		switch {
		case rawIn.Type == "context":
			msgIn = abm.Context{}
			if err = json.Unmarshal(data, &msgIn); err != nil {
				break
			}
			ctxt <- msgIn.(abm.Context)
		case rawIn.Type == "viewport":
			msgIn = render.Viewport{}
			if err = json.Unmarshal(data, &msgIn); err != nil {
				break
			}
			view <- msgIn.(render.Viewport)
		}
	}
}

func writer(ws *websocket.Conn, om <-chan goio.OutMsg) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case msg := <-om:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			jsonMsg, _ := json.MarshalIndent(msg, "", " ")
			ioutil.WriteFile("/tmp/outMsg", jsonMsg, 0644)
			// if err := ws.WriteJSON(msg); err != nil {
			// 	log.Println("writer: failed to WriteJSON:", err)
			// }
			_ = "breakpoint" // godebug
			time.Sleep(time.Millisecond * 100)
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Fatalln("writer: failed to WriteJSON:", err)
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Fatalln(err)
		}
		return
	}
	abm.InitModel(abm.DemoContext, abm.DemoEnvironment, om, view, phase)
	go writer(ws, om)
	reader(ws, ctxt, view)
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
