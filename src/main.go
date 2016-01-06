package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/gorilla/websocket"
)

// inMsg – typestring wrapper for generic *received* msg
type inMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// outMsh - typestring wrapper for *exported* msg
type outMsg struct {
	Type string `json:"type"`
	Data interface{}
}

const (
	// Time allowed to write the file to the client.
	writeWait = 20 * time.Millisecond

	// Time allowed to read the next pong message from the client.
	pongWait = 100 * time.Millisecond

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 1
	vpPopSize      = 0
	vsr            = d / 4
	γ              = 1.0
	cpplife        = -1
	vplife         = -1
	vpS            = 0.0
	vpA            = 1.0
	vτ             = quarterpi
	vκ             = 0.0
	v𝛔             = 0.0
	v𝛂             = 0.0
	cppS           = 0.01
	cppA           = 1.0
	cτ             = quarterpi
	sr             = d / 8
	randomAges     = false
	mf             = 0.5
	cφ             = 5
	cȣ             = 2
	cκ             = 1.0
	cβ             = 5
	vpAge          = false
	cppAge         = false
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	black = colour.Black
	white = colour.White

	e = abm.Environment{
		Bounds:         []float64{d, d},
		Dimensionality: dimensionality,
		BG:             black,
	}

	timeframe = abm.Timeframe{Turn: 0, Phase: 0, Action: 0}
	context   = abm.Context{
		e.Bounds,
		cppPopSize,
		vpPopSize,
		vpAge,
		vplife,
		vpS,
		vpA,
		vτ,
		vsr,
		γ,
		vκ,
		v𝛔,
		v𝛂,
		cppAge,
		cpplife,
		cppS,
		cppA,
		cτ,
		sr,
		randomAges,
		mf,
		cφ,
		cȣ,
		cκ,
		cβ,
	}
)

func reader(ws *websocket.Conn, context chan<- abm.Context, view chan<- render.Viewport) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var msgIn interface{}
	var rawIn JSONmsg
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
			context <- msgIn.(abm.Context)
		case rawIn.Type == "viewport":
			msgIn = render.Viewport{}
			if err = json.Unmarshal(data, &msgIn); err != nil {
				break
			}
			view <- msgIn.(render.Viewport)
		}
	}
}

func writer(ws *websocket.Conn, rm <-chan render.Msg) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case msg := <-rm:
			m, err := json.MarshalIndent(&msg, "", "  ")
			if err != nil {
				log.Fatalf("writer: failed when trying to marshal: %q", err)
			}
			err = ioutil.WriteFile("/tmp/dat1", m, 0644)
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteJSON(&m); err != nil {
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Fatalln(err)
		}
		return
	}

	go writer(ws)
	reader(ws)
}

func cppRBB(pop []abm.ColourPolymorphicPrey, queue chan<- render.AgentRender) {
	for {
		for i := 0; i < len(pop); i++ {
			c := &pop[i]
			𝚯 := calc.RandFloatIn(-c.Rτ, c.Rτ)
			c.Turn(𝚯)
			c.Move()
			c.Log()
			queue <- c.GetDrawInfo()
			time.Sleep(time.Second * 10)
		}
	}
}

func runningModel(m abm.Model, rc chan<- render.AgentRender, quit chan<- interface{}) {
	for {
		select {
		default:
			cppRBB(m.PopCPP, rc)
		}
	}
}

func initModel(context abm.Context) {
	simple := setModel(context)
	quit := make(chan interface{})
	rc := make(chan render.AgentRender)
	viz := make(chan render.Viewport)
	renderOut := make(chan render.Msg)

	go runningModel(simple, rc, quit)
	go visualiseModel(viz, rc, renderOut)
}

func setModel(c abm.Context) (m abm.Model) {
	m.PopCPP = abm.GeneratePopulation(cppPopSize, c)
	m.DefinitionCPP = []string{"mover"}
	m.Timeframe = timeframe
	m.Environment = e
	m.Context = c
	return
}

func visualiseModel(view <-chan render.Viewport, queue <-chan render.AgentRender, out chan<- render.Msg) {
	v := render.Viewport{"Viewport", 300, 200}
	tick := time.Tick(time.Second)
	var msg = render.Msg{"Render", nil, nil, e.BG}
	for {
		select {
		case job := <-queue:
			job = render.TranslateToViewport(job, v)
			msg.CPP = append(msg.CPP, job)
			out <- msg
			msg = render.Msg{"Render", nil, nil, e.BG} // reset msg contents
		case <-tick:
			// msg = render.Msg{} // reset msg contents
		case v = <-view:
		}
	}
}

func main() {

}
