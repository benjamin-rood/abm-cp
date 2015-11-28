package main

/*
MTime holds the model's representation of the time metrics.
Turn – The cycle length for all agents ∈ 𝐄 to perform 1 (and only 1) Action.
Phase – Division of a Turn, between agent sets, environmental effects/factors,
				and updates to populations and model conditions (via external).
				One Phase is complete when all members of a set have performed an Action
				or all requirements for the model's continuation have been fulfilled.
Action – An individual 'step' in the model. All Actions have a cost:
				the period (number of turns) before that specific Action can be
				performed again. For most actions this is zero.
				Some Actions could also *stop* any other behaviour by that agent
				for a period.
*/
type MTime struct {
	turn   int
	phase  int
	action int
}

const (
	x = iota
	y
	z
)

/*
ColRGB stores a standard 8-bit per channel Red Green Blue colour representation.
*/
type ColRGB struct {
	red   byte
	green byte
	blue  byte
}

/*
Environment specifies the boundary / dimensions of the working model. They
extend in both positive and negative directions, oriented at the center. Setting
any field (eg. zBounds) to zero will reduce the dimensionality of the model. For
most cases, a 2D environment will be sufficient.
*/
type Environment struct {
	xBounds float32
	yBounds float32
	zBounds float32
}

// Vector : Any sized dimension representation of an element of vector space.
type Vector []float32

// VectorCalc - Vector arithmetic, in (len(vector))dimension (i.e. any) space!

type colourPolymorhicPrey struct {
	pos       Vector
	movS      float32
	movA      float32
	heading   float32
	direction Vector
	lifetime  int32
	hunger    int32
	fertility int32 //	interval measurement between birth and sex
	gravid    bool  //	i.e. pregnant
	colour    ColRGB
}

type visualPredator struct {
	pos       Vector
	movS      float32
	movA      float32
	heading   float32 //	angle of direction relative to origin ∈ 𝐄
	direction Vector  //	unit vector for
	lifetime  int32   //	counter for number of turns agent exists in the model
	hunger    int32   //	measurement reflecting
	fertility int32   //	interval measurement between birth and sex
	gravid    bool    //	i.e. pregnant
	visRange  float32
	visAcuity float32
	imprint   ColRGB
}

// AgentActions interface for general agent behaviours
type AgentActions interface {
	UpdateSectorLocation()
	Turn(𝝧 float32)
	Move()
	Death()
}
