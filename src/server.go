package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/abm"
	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
	"github.com/gorilla/websocket"
)

type signal struct{}

// inMsg – typestring wrapper for generic *received* msg
type inMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// outMsh - typestring wrapper for *exported* msg
type outMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

const (
	// Time allowed to write the file to the client.
	writeWait = 50 * time.Millisecond

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	quarterpi      = 0.7853981633974483096156608458198757210492923498437764
	eigthpi        = 0.3926990816987241548078304229099378605246461749218882
	d              = 1.0
	dimensionality = 2
	cppPopSize     = 20
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
	ping     = signal{}
	om       = make(chan outMsg)
	phase    = make(chan signal)
	view     = make(chan render.Viewport)
	ctxt     = make(chan abm.Context)
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

	context = abm.Context{
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
	var rawIn inMsg
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

func writer(ws *websocket.Conn, om <-chan outMsg) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case msg := <-om:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteJSON(msg); err != nil {
				log.Fatalln("writer: failed to WriteJSON:", err)
			}
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
	initModel(context)
	go writer(ws, om)
	reader(ws, ctxt, view)
}

func cppRBB(pop []abm.ColourPolymorphicPrey, queue chan<- render.AgentRender) {
	for i := 0; i < len(pop); i++ {
		c := &pop[i]
		𝚯 := calc.RandFloatIn(-c.Rτ, c.Rτ)
		c.Turn(𝚯)
		c.Move()
		c.Log()
		queue <- c.GetDrawInfo()
		timeframe.Action++
	}
}

func runningModel(m abm.Model, rc chan<- render.AgentRender, quit <-chan signal, phase chan<- signal) {
	for {
		cppRBB(m.PopCPP, rc)
		timeframe.Phase++
		phase <- ping
		time.Sleep(time.Millisecond * 100)
	}
}

// hack for testing only
func initModel(context abm.Context) {
	simple := setModel(context)
	quit := make(chan signal)
	rc := make(chan render.AgentRender)
	go runningModel(simple, rc, quit, phase)
	go visualiseModel(view, rc, om, phase)
}

func setModel(c abm.Context) (m abm.Model) {
	m.PopCPP = abm.GeneratePopulation(cppPopSize, c)
	m.DefinitionCPP = []string{"mover"}
	m.Timeframe = timeframe
	m.Environment = e
	m.Context = c
	return
}

func visualiseModel(view <-chan render.Viewport, queue <-chan render.AgentRender, out chan<- outMsg, phase <-chan signal) {
	v := render.Viewport{300, 200}
	msg := outMsg{Type: "render", Data: nil}
	dl := render.DrawList{nil, nil, colour.RGB256{Red: 0, Green: 0, Blue: 0}}
	for {
		select {
		case job := <-queue:
			job.TranslateToViewport(v)
			switch job.Type {
			case "cpp":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			default:
				log.Fatalf("viz: failed to determine agent-render job type!")
			}
		case <-phase:
			msg.Data = dl
			out <- msg
			// reset msg contents
			msg = outMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{nil, nil, colour.RGB256{Red: 0, Green: 0, Blue: 0}}
		case v = <-view:
		}
	}
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
