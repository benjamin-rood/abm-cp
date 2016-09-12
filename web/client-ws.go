package web

import (
	"encoding/json"
	"time"

	"github.com/benjamin-rood/abm-cp/abm"
	"golang.org/x/net/websocket"
)

// SocketClient wraps the Model, bridging it to the WebSocket user connection.
type SocketClient struct {
	*websocket.Conn
	UUID string
	Name string
	*abm.Model
	active    bool
	timestamp time.Time
}

// NewSocketClient constructor
func NewSocketClient(ws *websocket.Conn, uuid string, params json.RawMessage) SocketClient {
	c := SocketClient{}
	c.Conn = ws
	c.UUID = uuid
	c.Name = "EMPTY"
	c.Model = abm.NewModel()
	c.active = true
	c.timestamp = time.Now()
	return c
}

// Monitor keeps the client's connection alive,
// and responds to any internal running model signaling
// – e.g. if there is a fault in the running abm,
// or if the population of the CP Prey agents reaches zero,
// then the model will invoke Kill() and Quit will close,
// which permits us to clean up and disconnect the SocketClient.
// Implements abm.Client interface.
func (c *SocketClient) Monitor(ch chan struct{}) {
	defer func() {
		c.active = false
		c.timestamp = time.Now() //	record when closed.
	}()
	for {
		select {
		case <-c.Quit: //	internal signal from Model
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

// Dead implements Client interface method for SocketClient
func (c *SocketClient) Dead() bool {
	return c.active
}

// TimeStamp implement Client interface method for SocketClient
func (c *SocketClient) TimeStamp() time.Time {
	return c.timestamp
}
