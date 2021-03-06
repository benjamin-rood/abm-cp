package abm

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path"
  "path/filepath"
  "time"
)

// Data Logging process local to the model instance.
func (m *Model) log(ec chan<- error) {
  time.Sleep(pause)

  signature := "LOG_" + m.SessionIdentifier
  turnEnd, clash := m.turnSync.Register(signature)
  if clash {
    errStr := "Clash when registering Model: " + m.SessionIdentifier + " log: for sync with m.turnSync"
    ec <- errors.New(errStr)
    return
  }

  defer func() {
    m.turnSync.Deregister(signature)
    // Need to wipe the agent records too? -yes, probably... but should be in Stop()?
  }()

  if m.UseCustomLogPath {
    m.LogPath = path.Join(os.Getenv("HOME")+os.Getenv("HOMEPATH"), m.CustomLogPath, abmlogPath, m.SessionIdentifier, m.timestamp)
  } else {
    m.LogPath = path.Join(os.Getenv("HOME")+os.Getenv("HOMEPATH"), abmlogPath, m.SessionIdentifier, m.timestamp)
  }

  for {
    select {
    case <-m.halt: // RUN halted as channel closed – therefore we end LOG.
      // clean up?
      return
    case <-turnEnd:
      func() {
        cpr := m.cpPreyRecordCopy()
        vpr := m.vpRecordCopy()
        go func(record map[string]ColourPolymorphicPrey, errCh chan<- error) {
          // write map as json to file.
          tc := fmt.Sprintf("%08v", m.Turn)
          dir := m.LogPath
          path := dir + string(filepath.Separator) + tc + "_cpPrey_pop_record.dat"

          msg, err := json.MarshalIndent(record, "", "  ")
          if err != nil {
            log.Printf("model: logging: json.Marshal failed, error: %v\n source: %s : %s : %v\n", err, m.SessionIdentifier, m.timestamp, m.Turn)
            errCh <- err
            return
          }

          err = os.MkdirAll(dir, 0777)
          if err != nil {
            errCh <- err
            return
          }
          err = ioutil.WriteFile(path, msg, 0777)
          if err != nil {
            errCh <- err
            return
          }
        }(cpr, ec)
        go func(record map[string]VisualPredator, errCh chan<- error) {
          // write map as json to file.
        }(vpr, ec)
      }()
    }
  }
}
