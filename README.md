[![Go Report Card](https://goreportcard.com/badge/github.com/benjamin-rood/abm-cp)](https://goreportcard.com/report/github.com/benjamin-rood/abm-cp)

###Note that all code is alpha until version 1.0

###Current Version: 0.4.6

GIF Preview: https://gfycat.com/NippyTidyCassowary

YouTube Preview: https://www.youtube.com/watch?v=FJfskjJOQKc

Software for computing a Predator-Prey Agent Based Model of prey colour polymorphism (CP) in Go. 

Written specifically to assist the research by Dr. James Dale in evaluating hypotheses for the evolutionary maintenance of extreme colour polymorphism.

[Initial Academic Report](https://goo.gl/GjQqv6)

##Goal

Server-side computation, client-side GUI controls and visualisation of the running ABM in a web browser.

Generalise system so any `abm` package with local model conditions, agent design, etc can be hot-swapped in (linked via a command line flag) and run.

###How to use:

Install Go from [here](https://golang.org/dl/).

Download my repo: `go get -u github.com/benjamin-rood/abm-cp`

Download external dependencies: 

`go get -u golang.org/x/net/websocket`
`go get -u github.com/benjamin-rood/gobr`
`go get -u github.com/pquerna/ffjson`
`go get -u github.com/davecgh/go-spew/spew`
`go get -u github.com/spf13/cobra`

*(All dependecies will be vendored into the `abm-cp` package from v1.0.0 onwards)*

Change current directory to `$GOPATH/src/github.com/benjamin-rood/abm-cp` and run `go install`.
From there, calling `abm-cp run` from the shell prompt will start the websocket server.

Point web browser at `localhost:8000`

Current version only tested on Safari on OS X.



## Roadmap

### 0.1.0
Base requirements completed.

* A simple interface for the CP Prey ABM, with the visualisation on the left, with conditionsual parameter controls on the right. :white_check_mark:

* Use P5js for render viewport. :white_check_mark:

* Browser recieves drawing instructions from ABM which get loaded into the P5 instance's draw loop. :white_check_mark:

* Responsive design, visualisation (render viewport) scales to the dimensions of available browser real estate. :white_check_mark:

* Server handles concurrent bi-directional connections for concurrent ABM modelling individual to each user, with data sent between client and server using Web Sockets. :white_check_mark:

* Server cleans up after user disconnections. :white_check_mark:

* Server receives new submitted conditionsual parameters as signal to start/restart ABM. :white_check_mark:

* Serialisation of data messages using JSON (prohibatively slow for anything approaching 10,000 agents).  :white_check_mark:

* CP Prey agents implementation:
	 * Rule-Based-Behaviour. :white_check_mark:
	 * Asexual Reproduction. :white_check_mark:
	 * Mutation (colouration). :white_check_mark:

* Visual Predator implementation:
	* Rule-Based-Behaviour. :white_check_mark:
	* Exhaustive Prey Search (very slow). :white_check_mark:
	* Colour Imprinting (needs tweaking, no baseline yet established). :white_check_mark:

* Simple concurrent execution of Predator/Prey Action. :white_check_mark:

* Unit tests for `geometry`, `calc`, `render` packages. :white_check_mark:

### 0.2.0

* Dispatch errors along channels, use external handlers who receive the errors and process/print them. :white_check_mark:

* Essential unit tests for `abm` package :white_check_mark:

* Show live population and timeline statistics inside P5js viewport. :white_check_mark:
 
* Visual Predator implemenation:
	* Find baseline params for Colour Imprinting. :white_check_mark:
	* Adaptation in response to hunger. :white_check_mark:
	* Starvation ⟹ Death. :white_check_mark:
	* Sexual Reproduction. :white_check_mark:

* General modelling and interactions between agent types in place, with all baseline parameters set for end-use. :white_check_mark:

### 0.3.0

* Toggle Visualisation on/off :white_check_mark:

* Record data of a running model to log files as JSON for statistical analysis as a concurrent process, independent of Visualisation or model computation. Toggleable. :white_check_mark:

* User determined frequency of data logging on client side. :white_check_mark:

* Better model for Predator specialisation through colour-imprinting which directly gives search/attack advantage, rather than being decided randomly. :white_check_mark:

* Complete tests for `abm` package :white_check_mark:

### 0.4.0

* Use `ffjson`–generated custom Marshal/Unmarshal JSON methods for ~2X speedup when serialising render messages to client  :white_check_mark:

* Better Prey Search using 2d dimensional search trees.

* Browser-side input validation.

* Beging switch to `spf13/cobra` CLI system. :white_check_mark:

* Use k-dimensional tree for spatial partitioning of model environment, permitting optimal search. 
(General implementation in Go already done using trees connected to pointers to elements in slices)

* Web-side input validation for web contextual parameters

### 0.4.0

* Have complete control over ABM computation, logging, visualisation from command-line, rather than just starting up a web server and controlling through the the browser.

* Use uncompressed JSON-formatted logging for debug only.

* Switch to a compressed binary encoding for log files – or try FlatBuffers?

### 0.4.0

* Switch data serialisation to Protocol Buffers (protobuf) ~10X speedup. Marshalling drawing instructions to JSON is currently the single most expensive action!

###1.0.0 – Late 2016?

**Use Amazon Web Services and switch to a model of cloud (distributed) computation and storage** for all log files, thus entirely taking the burden off the user for all hardware costs in the modelling. Whilst the need for CPU and memory optimisation along with data compression over the wire remains essential, the scale of the model environment and populations could become entirely unrestricted.

* Switch all public html file to templated/generated ones based on conditions parameters etc.

* Switch to `gopherjs` for all front-end code?

* Allow end-user to switch between different browser layouts: Visualisation only, Standard and Statistical? $\Leftarrow$ Could use *Jupyter* to present graphing in browser?

* Start ABM computation remotely and keep running after disconnection? i.e. start the model running, leave it, reconnect based on session UUID at a later tim to check up or review results.

* Batch processing.

* Email user when model session finishes.

* Enable use in a distributed environment.

* Complete testing suite including integration tests.

* Allow hot-swapping of different `abm` variant packages.

* Store modelling sessions to server database along with statistical data for later retrieval.

* Fluid ABM timescale controls?  Doable, but probably not without switching to gopherjs so I can integrate it all within the same codebase.

* Optional recording of Visualisation to SVG frame sequence using `ajstarks/svgo` package.
