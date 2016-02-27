
##Implementation Overview

![intial whiteboard](http://i.imgur.com/VZdmZXf.jpg)

<br>
After consideration and discussion with Assoc. Prof. Scogings, the following conclusions emerged:

* The relationship between the PREDATOR and PREY agents can be modeled exclusively in two dimensions.

* To increase flexibility, and allow for spatial interactions with greater nuance, agents should operate in a discrete $\mathbb{R}^{2}$ vector space, rather than a fixed integer field using matrices to determine agent positioning.

* Focus on concurrent ABM modelling/computation (rule-based agent behaviour),  visualisation, and data logging for statistical analysis.

* The surface area to parallelise computation of the domain is very narrow, as agent actions are largely interdependent of one another at each moment in time. Instead, implementation requires *concurrent software design* from the beginning.

* Visualisation must be in "Live" viewport

* Must be GUI application — interactive controls to sett constant value parameters which are received and then used to either update or (re)start model.

* Software must be designed for modularity: need to support extensibility of agent behaviours and parameters per J. Dale specifications to ensure software correct meets research needs.

* The best endpoint for the user is a Web App:
	* Platform agnostic
	* Every possible user has a web browser on their device
	* GUI stable and low-overhead, programmable by CSS, HTML, JS
	* Allows for local and remote execution — both cases, same UI & UX

I made the decision to develop the application software in [**Go**](golang.org), a language designed by Rob Pike, Ken Thompson, and Robert Griesemer at Google[^golang.org - history] and released to the public in 2009. It was designed for composing concurrent web/cloud backend software using concepts from *Communicating Sequential Processes* by Tony Hoare.[^golang.org - ancestors] [^Hoare]
The front-end would naturally be presented as a web page written in HTML, CSS, and Javascript. In particular, after an exhaustive search, the [**P5js**](http://p5js.org/) library was chosen for its simplicity and emphasis on 2D drawing, in order to render the visualisation viewport inside the browser window.




[^golang.org - history]: <sub> https://golang.org/doc/faq#history

[^golang.org -ancestors]: <sub> https://golang.org/doc/faq#ancestors

[^Hoare]: <sub> Hoare, C. A. R.(Charles Anthony Richard) (1985). *Communicating Sequential Processes*.  Englewood Cliffs, N.J. : Prentice/Hall International