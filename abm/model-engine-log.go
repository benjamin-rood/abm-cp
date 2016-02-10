package abm

import (
	"bytes"
	"encoding/json"
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
	fmt.Println("starting logging...")
	time.Sleep(pause)
	_ = "breakpoint" // godebug
	var turn chan struct{}
	var clash bool
	var signature string
	defer func() {
		m.turnSignal.Deregister(signature)
		// Need to wipe the agent records too? -yes, probably.
	}()

Registration: // register to receive from m.turnSignal - loop until no clash with existing entry in map.
	signature = uuid()
	turn, clash = m.turnSignal.Register(signature)
	if clash {
		goto Registration
	}

	if m.UseCustomLogPath {
		m.LogPath = path.Join(os.Getenv("HOME")+os.Getenv("HOMEPATH"), m.CustomLogPath, abmlogPath, m.SessionIdentifier, m.timestamp)
	}

	for {
		select {
		case <-m.rc: // run finished as rc channel closed!
			time.Sleep(time.Second)
			// clean up?
			return
		case <-turn:
			func() {
				reccpp := m.copyCppRecord()
				recvp := m.copyVpRecord()
				go func(rc map[string]ColourPolymorphicPrey, errCh chan<- error) {
					path := m.LogPath + string(filepath.Separator) + "0" + "_vp_pop_record.dat"

					msg, err := json.MarshalIndent(rc, "", "  ")
					if err != nil {
						log.Printf("model: logging: json.Marshal failed, error: %v\n source: %s : %s : %v\n", err, m.SessionIdentifier, m.timestamp, m.Turn)
						errCh <- err
						return
					}

					var buff []byte
					out := bytes.NewBuffer(buff)
					out.Write(msg)
					output := make([]byte, 1024*10)
					n, rerr := out.Read(output)
					if n == 0 || rerr != nil {
						fmt.Println("n:", n, "rerr:", rerr.Error())
						errCh <- err
						return
					}

					err = os.MkdirAll(m.LogPath, 0777)
					if err != nil {
						log.Println(err)
						errCh <- err
						return
					}
					err = ioutil.WriteFile(path, output, 0777)
					if err != nil {
						// fmt.Println(err)
						errCh <- err
						return
					}
				}(reccpp, ec)
				go func(rv map[string]VisualPredator) {
					// write map as json to file.
				}(recvp)
			}()
		}
	}
}