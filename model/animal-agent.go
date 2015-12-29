package model

import "github.com/benjamin-rood/abm-colour-polymorphism/geometry"

// Interfaces which define an animal type agent.

// Mover defines an agent whose position is non-static
type Mover interface {
	Turn(𝚯 float64)
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
	MateSearch([]Breeder) (bool, *Breeder)
	Copulation(*Breeder) bool
	Birth() []Breeder
}

// Mortal defines an agent which ages and dies.
type Mortal interface {
	Age()
	Death()
}
