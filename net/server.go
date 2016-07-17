package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benjamin-rood/abm-cp/abm"
	"github.com/benjamin-rood/gobr"
	"golang.org/x/net/websocket"
)

const (
	sweepFreq   = time.Minute
	deathPeriod = time.Hour * 24
)

// global mutable index of current abm-cp users.
var socketUsers = make(map[string]abm.WebScktClient)

func clientUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func sweepSocketWebScktClients() {
	sweeper := time.NewTicker(sweepFreq)
	select {
	case <-sweeper.C:
		for uid, client := range socketUsers {
			if client.Dead {
				delete(socketUsers, uid)
				continue
			}
			if time.Since(client.Stamp) >= deathPeriod {
				delete(socketUsers, uid)
			}
		}
	}
}

func wsSession(ws *websocket.Conn) {
	uuid := clientUUID()
	log.Println("wsSession uuid:", uuid)
	c := abm.NewWebScktClient(ws, uuid)
	socketUsers[uuid] = c
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Println("wsSession exit failed to close Conn!", err)
		}
		delete(socketUsers, uuid)
	}()
	wsCh := make(chan struct{})
	go c.ErrPrinter()
	go c.Controller()
	go wsReader(ws, c.Im, wsCh)
	go wsWriter(ws, c.Om, wsCh)
	c.Monitor(wsCh) //	keep session alive
}

// TODO: don't pass the raw CONN, instead, pass the *client* ?
func wsReader(ws *websocket.Conn, in chan<- gobr.InMsg, quit chan struct{}) {
	_ = "breakpoint" // godebug
	defer func() {
		//	clean up
	}()
	for {
		msg := gobr.InMsg{}
		select {
		case <-quit:
			return
		default:
			err := websocket.JSON.Receive(ws, &msg)
			log.Printf("received JSON msg: %s\n", msg)
			if err != nil {
				log.Println("error: wsReader:", err)
				log.Println("Disconnected User.")
				close(quit)
				return
			}
			in <- msg
		}
	}
}

func wsWriter(ws *websocket.Conn, out <-chan gobr.OutMsg, quit <-chan struct{}) {
	defer func() {
		// clean up
	}()

	for {
		select {
		case <-quit:
			return
		case msg := <-out:
			websocket.JSON.Send(ws, msg)
			//spew.Dump(msg)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocket.Handler(wsSession))
	http.ListenAndServe(":9999", nil)
}
