package web

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benjamin-rood/abm-cp/abm"
	"github.com/benjamin-rood/gobr"
	"golang.org/x/net/websocket"
)

const (
	sweepFreq   = time.Minute * 5
	deathPeriod = time.Hour * 24
)

// Users the global mutable index of current abm-cp users.
var Users = make(map[string]abm.Client)

func clientUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func sweepClients() {
	sweeper := time.NewTicker(sweepFreq)
	select {
	case <-sweeper.C:
		for uid, client := range Users {
			if client.Dead() {
				delete(Users, uid)
				continue
			}
			if time.Since(client.TimeStamp()) >= deathPeriod {
				delete(Users, uid)
			}
		}
	}
}

func wsSession(ws *websocket.Conn) {
	uuid := clientUUID()
	log.Println("wsSession uuid:", uuid)
	paramsJSON, _ := json.MarshalIndent(abm.DefaultConditionParams, "", " ")
	c := NewSocketClient(ws, uuid, json.RawMessage(paramsJSON))
	Users[uuid] = &c
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Println("wsSession exit failed to close Conn!", err)
		}
		delete(Users, uuid)
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
		}
	}
}

// WsServer is the process launched by `abm-cp` program by default `run` command
// soon will be
func WsServer(port string, ws string) {
	http.Handle("/", http.FileServer(http.Dir("./web/public")))
	http.Handle(`/`+ws, websocket.Handler(wsSession))
	http.ListenAndServe(`:`+port, nil) // need (channel-based?) way of closing this.
}
