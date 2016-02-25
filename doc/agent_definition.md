
## Agent Classification & Definition ##

Given the context of the ABM, it should be useful to adopt some of the language of biological terminology to conceive of a particular type of agent and its behaviour when writing the code necessary to implement the algorithms defining the _rule-based behaviour_ of each agent type.

Therefore, let the classification of an agent type be (broadly) analogous to the taxonomy of biological organisms and their general characteristics, s.t. we may say Agent Type $\equiv$ Species delimination.

Further, that the implementation of that classification and definition of an agent be determined by which _interfaces_ it has implemented.

In fact, it turns out that Go's interfaces, and the language's emphasis on _components_ and the "contract" of expectation of behaviour via interfaces is what actually suggests such an approach.

	// Interfaces which define an animal type agent.
	
	// Mover defines an agent whose position is non-static
	type Mover interface {
		Turn(ϴ float64)
		Move(target geometry.Vector)
	}
	
	// Hunter defines an agent which looks for sustinence by going after prey agents.
	type Hunter interface {
		PreySearch() (bool, Hunted)
		Attack(Hunted) bool
		Eat(Hunted) bool
	}
	
	// Forager defines an agent which looks for sustinence by searching its environment – although it can be carnivorous, it does not Hunt for live prey.
	type Forager interface {
		FoodSearch()
		Eat()
	}
	
	// Hunted defines an agent which must avoid Hunters!
	type Hunted interface {
		Evade(Hunter) //	e.g. dodging!
		Hide()        //	e.g. uses camoflage
	}
	
	// Defender defines an agent which actively repels attacks
	type Defender interface {
		Block() bool
		Counter()
	}
	
	// Breeder defines an agent which breeds sexually with other agents of the same type.
	type Breeder interface {
		MateSearch([]Breeder) *Breeder
		Copulation(*Breeder) bool
		Birth() []Breeder
	}
	
	// Mortal defines an agent which ages and dies.
	type Mortal interface {
		Age()
		Death()
	}