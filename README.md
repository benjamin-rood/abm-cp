
###Current Version: 0.1.5

![abm preview](https://giant.gfycat.com/AggressiveGregariousAidi.gif)

##Context
Software for computing a Predator-Prey Agent Based Model of prey colour polymorphism (CP) in Go. 

Written specifically to assist the research by Dr. James Dale in evaluating hypotheses for the evolutionary maintenance of extreme colour polymorphism.

##Goal

Server-side computation, client-side GUI controls and visualisation of the running ABM in a web browser.

Generalise system so any `abm` package with local model context, agent design, etc can be hot-swapped in (linked via a command line flag) and run.

###How to use:
Install Go from [here](https://golang.org/dl/).
Download this repository:  run `go get -u github.com/benjamin-rood/abm-colour-polymorphism`
Change current directory to `$GOPATH/src/github.com/benjamin-rood/abm-colour-polymorphism` and run `cd net && go build && ./net`
Point web browser at `localhost:8080`

Only tested so far on Safari (OS X) and Chrome (OS X, Windows 8.1).



## Roadmap

### 0.1.0
Base requirements completed.

* A simple interface for the CP Prey ABM, with the visualisation on the left, with contextual parameter controls on the right. :white_check_mark:

* Use P5js for render viewport. :white_check_mark:

* Browser recieves drawing instructions from ABM which get loaded into the P5 instance's draw loop. :white_check_mark:

* Responsive design, visualisation (render viewport) scales to the dimensions of available browser real estate. :white_check_mark:

* Server handles concurrent bi-directional connections for concurrent ABM modelling individual to each user, with data sent between client and server using Web Sockets. :white_check_mark:

* Server cleans up after user disconnections. :white_check_mark:

* Server receives new submitted contextual parameters as signal to start/restart ABM. :white_check_mark:

* Serialisation of data messages using JSON (prohibatively slow for anything approaching 10,000 agents).  :white_check_mark:

* CP Prey agents implementation:
	 * Rule-Based-Behaviour. :white_check_mark:
	 * Asexual Reproduction. :white_check_mark:
	 * Mutation (colouration). :white_check_mark:

* Visual Predator implementation:
	* Rule-Based-Behaviour. :white_check_mark:
	* Exhaustive Prey Search (very slow). :white_check_mark:
	* Colour Imprinting (needs tweaking, no baseline yet established). :white_check_mark:

* Simple concurrent execution of Predator/Prey RBB. :white_check_mark:

* Unit tests for `geometry`, `calc`, `render` packages. :white_check_mark:

### 0.2.0

* Essential unit tests for `abm` package :white_check_mark:

* Show live population (and timeline) statistics inside P5js viewport.

* Client-side input validation.

* Switch data serialisation to Protocol Buffers (protobuf) – marshalling drawing instructions to JSON is currently the single most expensive action!
 
* Visual Predator implemenation:
	* Find baseline params for Colour Imprinting. :white_check_mark:
	* Adaptation in response to hunger. :white_check_mark:
	* Starvation ⟹ Death. :white_check_mark:
	* Sexual Reproduction. :white_check_mark:

* Expected modelling entirely in place, with all baseline parameters in place. :white_check_mark:

### 0.3.0

* Statistical logging to file using IOWriter.

* Better Prey Search (using grid system).

* Complete tests for `abm` package


### 1.0.0

* Use *k-dimensional tree* for spatial partitioning of model environment.

* Live statistical graphing widgets

* Allow end-user to switch between different browser layouts: Visualisation only, Standard, and Statistical.

* Control ABM computation from command-line.

* Enable use in a distributed environment.

* Complete testing suite.

* Allow hot-swapping of different `abm` packages.

* Store modelling sessions to server database along with statistical data for later retrieval.

*  Offline ABM computation (start the model running, leave it, reconnect, see what it's up to).

* Fluid ABM timescale controls.

* Optional recording of Visualisation to SVG frame sequence. 