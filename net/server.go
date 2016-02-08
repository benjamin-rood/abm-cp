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

// Client defines the owner of the Model.
type Client struct {
	*websocket.Conn
	UUID string
	Name string
	*abm.Model
	Active bool
	Stamp  time.Time
	Quit   chan struct{}
}

// NewClient constructs an initialised Client session.
func NewClient(ws *websocket.Conn, uuid string) (c Client) {
	c.Conn = ws
	c.UUID = uuid
	c.Name = name()
	c.Model = abm.NewModel()
	c.Active = true
	c.Stamp = time.Now()
	c.Quit = make(chan struct{})
	return
}

// Monitor keeps the client's connection alive,
// and responds to any internal running model signaling
// – e.g. if there is a fault in the running abm,
// or if the population of the CP Prey agents reaches zero,
// then the model will invoke Kill() and Quit will close,
// which permits us to clean up and disconnect the Client.
func (c *Client) Monitor(ch chan struct{}) {
	defer func() {
		c.Active = false
		c.Stamp = time.Now()
	}()
	for {
		select {
		case <-c.Quit: //	internal signal from client.
			close(ch) //	exit websocket connection.
			// send final statistics
			// clean up
			return
		case <-ch:
			c.Suspend() //	websocket connection dead, suspend model operation.
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func name() string {
	return "01b"
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func networkError(err error, c chan struct{}) {
	log.Println(err)
	close(c)
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

const (
	sweepFreq   = time.Minute
	deathPeriod = time.Hour * 24
)

// global mutable index of current users.
var socketUsers = make(map[string]Client)

func sweepSocketClients() {
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

// TODO: too hackneyed?
func wsSession(ws *websocket.Conn) {
	uuid := uuid()
	log.Println("wsSession uuid:", uuid)
	c := NewClient(ws, uuid)
	socketUsers[uuid] = c
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			log.Println("wsSession exit failed to close Conn!", err)
		}
		delete(socketUsers, uuid)
	}()
	wsCh := make(chan struct{})
	go wsReader(ws, c.Im, wsCh)
	go wsWriter(ws, c.Om, wsCh)
	go c.Controller()
	c.Monitor(wsCh) //	keep alive
}

// TODO: don't pass the raw CONN! Instead, pass the *client*
func wsReader(ws *websocket.Conn, in chan<- gobr.InMsg, quit chan struct{}) {
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
	http.ListenAndServe(":8080", nil)
}
