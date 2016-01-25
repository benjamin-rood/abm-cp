package main

import (
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"

	"golang.org/x/net/websocket"
)

// Client defines the owner of the Model.
type Client struct {
	*websocket.Conn
	UID string
	abm.Model
	Active bool
	Stamp  time.Time
	Quit   chan struct{}
}

// NewClient constructs an initialised Client session.
func NewClient(ws *websocket.Conn, uid string) (c Client) {
	c.Conn = ws
	c.UID = uid
	c.Model = abm.NewModel()
	c.Active = true
	c.Stamp = time.Now()
	c.Quit = make(chan struct{})
	return
}

func (c *Client) Monitor(ch chan struct{}) {
	defer func() {
		c.Active = false
		c.Stamp = time.Now()
	}()
	for {
		select {
		case <-c.Quit: //	internal kill signal from client.
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
