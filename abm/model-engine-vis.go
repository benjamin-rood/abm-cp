package abm

import (
	"errors"
	"fmt"

	"github.com/benjamin-rood/abm-cp/render"
	"github.com/benjamin-rood/gobr"
)

// Visualisation process local to the model instance.
func (m *Model) vis(ec chan<- error) {
	signature := "VIS_" + m.SessionIdentifier
	turn, clash := m.turnSignal.Register(signature)
	if clash {
		errStr := "Clash when registering Model: " + m.SessionIdentifier + " vis: for sync with m.turnSignal"
		ec <- errors.New(errStr)
		return
	}
	defer m.turnSignal.Deregister(signature)

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
			// fmt.Println("VIS_ received on render channel:")
			// spew.Dump(job)
			switch job.Type {
			case "cpPrey":
				dl.CPP = append(dl.CPP, job)
			case "vp":
				dl.VP = append(dl.VP, job)
			}
		case <-turn:
			dl.CppPop = fmt.Sprintf("cpPrey %d", len(m.PopCPP))
			dl.VpPop = fmt.Sprintf("vp  %d", len(m.PopVP))
			dl.TurnCount = fmt.Sprintf("%08d", m.Turn)
			msg.Data = dl
			// fmt.Println("VIS_ sending out on Om channel:")
			// spew.Dump(msg)
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
			return
		}
	}
}
