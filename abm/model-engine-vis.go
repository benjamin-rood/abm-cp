package abm

import (
	"fmt"
	"time"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

// Visualisation process local to the model instance.
func (m *Model) vis(ec chan<- error) {
	var turn chan struct{}
	var clash bool
	var signature string
	defer m.turnSignal.Deregister(signature)

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

	msg := gobr.OutMsg{Type: "render", Data: nil}
	bg := m.BG.To256()
	dl := render.DrawList{
		CPP:       nil,
		VP:        nil,
		BG:        bg,
		CppPop:    "0",
		VpPop:     "0",
		TurnCount: "0",
	}
	for {
		select {
		case job := <-m.render:
			switch job.Type {
			case "cpp":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			}
		case <-turn:
			dl.CppPop = fmt.Sprintf("cpp %d", len(m.PopCPP))
			dl.VpPop = fmt.Sprintf("vp  %d", len(m.PopVP))
			dl.TurnCount = fmt.Sprintf("%08d", m.Turn)
			msg.Data = dl
			m.Om <- msg
			// reset msg contents
			msg = gobr.OutMsg{Type: "render", Data: nil}
			//	reset draw instructions
			dl = render.DrawList{ //	reset
				CPP:       nil,
				VP:        nil,
				BG:        bg,
				CppPop:    "0",
				VpPop:     "0",
				TurnCount: "0",
			}
		case <-m.rc: //	run channel closed!
			time.Sleep(time.Second)
			return
		}
	}
}
