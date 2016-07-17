package abm

import (
	"time"

	"golang.org/x/net/websocket"
)

// WebScktClient wraps the Model, bridging it to the WebSocket user connection.
type WebScktClient struct {
	*websocket.Conn
	UUID string
	Name string
	*Model
	Active bool
	Dead   bool
	Stamp  time.Time
}

// NewWebScktClient constructor
func NewWebScktClient(ws *websocket.Conn, uuid string) WebScktClient {
	c := WebScktClient{}
	c.Conn = ws
	c.UUID = uuid
	c.Name = "EMPTY"
	c.Model = NewModel()
	c.Active = true
	c.Dead = false
	c.Stamp = time.Now()
	return c
}

// Monitor keeps the client's connection alive,
// and responds to any internal running model signaling
// – e.g. if there is a fault in the running abm,
// or if the population of the CP Prey agents reaches zero,
// then the model will invoke Kill() and Quit will close,
// which permits us to clean up and disconnect the WebScktClient.
func (c *WebScktClient) Monitor(ch chan struct{}) {
	defer func() {
		c.Active = false
		c.Stamp = time.Now()
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
