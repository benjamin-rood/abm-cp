package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

type positiveValue struct {
	set   bool
	val   int
	label string
}

func (pv *positiveValue) Set(x string) (err error) {
	if pv == nil {
		return errors.New("positiveValue: <nil>")
	}
	pv.val, err = strconv.Atoi(x)
	if err != nil {
		return err
	}
	pv.set = true
	return positiveValueErr(*pv)
}

func (pv *positiveValue) String() string {
	return fmt.Sprintf("%d", pv.val)
}

func (pv *positiveValue) Get() interface{} {
	return pv
}

type ratioParamErr struct {
	Value string
	Kind  string
}

func (e *ratioParamErr) Error() string {
	errStr := "positiveValue: (" + e.Value + ") is an ipvalid parameter value.\n"
	switch e.Kind {
	case "render":
		errStr += (visFreqFlagUsage + visFreqFlagDescription)
	case "log":
		errStr += logFreqFlagUsage
	}
	return errStr
}

func positiveValueErr(pv positiveValue) error {
	switch {
	case (pv.val < 0):
		return &ratioParamErr{"number %v", pv.label}
	}
	return nil
}

// -- bool Value
type boolValue struct {
	set   bool
	val   bool
	label string
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	b.val = v
	b.set = true
	return err
}

func (b *boolValue) Get() interface{} { return b }

func (b *boolValue) String() string { return fmt.Sprintf("%v", *b) }

func (b *boolValue) IsBoolFlag() bool { return true }

type command struct {
	Description string
	flag.Value
}

const (
	wsFlagLabel              = `start websocket server`
	portFlagLabel            = `defines webserver port`
	modelConditionsFlagLabel = `input model ConditionParams from JSON`
	visFreqFlagLabel         = `set render frequency`
	logFreqFlagLabel         = `set logging frequency`
	logDetailFlagLabel       = `set logging detail`
	logOutputPathFlagLabel   = `set logging output path`

	wsFlags              = `-w | --web`  // boolean flag, default = false
	portFlags            = `-p | --port` // default = 80
	modelConditionsFlags = `-m | -c | --conditions <JSON>`
	visFreqFlags         = `-v | --vis <frequency>`
	logFreqFlags         = `-l | --log <frequency>`
	logDetailFlags       = `-ld | -- logdetail <integer>`
	logOutputPathFlags   = `-lo | --output <directory path>`

	wsFlagDescription      = wsFlags + `: Boolean flag to create a Web GUI for the Agent-Based Model session via WebSockets. On successful initialisation, program will return a URL to the unique session. (Default = false)`
	visFreqFlagDescription = visFreqFlags + `: Optional flag followed by a value T, defines visualisation/rendering ratio T:1 where T is the number of ABM "turns" which are computed for every 1 render of the Model State, e.g. (T=4) = 4:1 = 1 rendered frame per 4 turns. (Default (T=1) = 1:1 ratio)`
	logFreqFlagDescription = logFreqFlags + `: Optional flag followed by a value T, defines the logging ratio T:1 where T is the number of ABM turns which are computed for every log record of the state of both agent populations, e.g. (T=10) = 10:1 = 1 complete log record of all agents of each type every 10 turns that elapse.
	(Default (T=1) = 1:1 ratio)`

	wsFlagUsage      = `usage: [` + wsFlags + `]` + ` (` + wsFlagLabel + `)`
	portFlagUsage    = `usage: [` + portFlags + `]` + ` (` + portFlagLabel + `)`
	visFreqFlagUsage = `usage: [` + visFreqFlags + `]` + ` (` + visFreqFlagLabel + `: where the interval parameter must either be a positive integer or zero (signifying that rendering is disabled.)`
	logFreqFlagUsage = `usage: [` + logFreqFlags + `]` + ` (` + logFreqFlagLabel + `)`

	usageString = `usage: abm-cp [` + wsFlags + `] [` + portFlags + `] [` + modelConditionsFlags + `] [` + visFreqFlags + `] [` + logFreqFlags + `] [` + logDetailFlags + `] [` + logOutputPathFlags + `] [-h | --help] `
)

func main() {
	// modellingParameters := abm.DefaultConditionParams
	v := command{Description: visFreqFlagDescription, Value: &positiveValue{false, 1, visFreqFlagLabel}}
	w := command{Description: "", Value: &boolValue{false, false, wsFlagLabel}}
	l := command{Description: logFreqFlagDescription, Value: &positiveValue{false, 1, logFreqFlagLabel}}

	flag.Var(v.Value, "vis", visFreqFlagUsage)
	flag.Var(v.Value, "v", visFreqFlagUsage)
	flag.Var(l.Value, "log", logFreqFlagUsage)
	flag.Var(w.Value, "web", wsFlagUsage)
	flag.Var(w.Value, "w", wsFlagUsage)

	if flag.NFlag() == 0 {
		log.Println(usageString)
		os.Exit(0)
	}
}
