package abm

import (
  "errors"
  "fmt"

  "github.com/benjamin-rood/abm-cp/render"
  "github.com/benjamin-rood/gobr"
)

// Visualisation process local to the model instance, orchestrated with the RUN process
func (m *Model) vis(ec chan<- error) {
  signature := "VIS_" + m.SessionIdentifier
  turnEnd, clash := m.turnSync.Register(signature)
  if clash {
    errStr := "Clash when registering Model: " + m.SessionIdentifier + " vis: for sync with m.turnSync"
    ec <- errors.New(errStr)
    return
  }
  defer m.turnSync.Deregister(signature)

  msg := gobr.OutMsg{Type: "render", Data: nil}
  bg := m.BG.To256()
  dl := render.DrawList{
    CPP:       nil,
    VP:        nil,
    BG:        bg,
    CpPreyPop: "0",
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
    case <-turnEnd:
      dl.CpPreyPop = fmt.Sprintf("cpPrey %d", len(m.popCpPrey))
      dl.VpPop = fmt.Sprintf("vp  %d", len(m.popVisualPredator))
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
        CpPreyPop: "0",
        VpPop:     "0",
        TurnCount: "0",
      }
    case <-m.halt: //	RUN stopped as channel closed – therefore we end VIS.
      return
    }
  }
}
