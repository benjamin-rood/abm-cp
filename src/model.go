package main

// ColRGB is a standard 8-bit per channel colour representation.
type ColRGB struct {
	r byte
	g byte
	b byte
}

// Environment specifies the boundary / dimensions of the working model. They extend in both positive and negative directions, oriented at the center. Setting any field (eg. zBounds) to zero will reduce the dimensionality of the model. For most cases, a 2D environment will be sufficient.
type Environment struct {
	xBounds float32
	yBounds float32
	zBounds float32
}

// Vec2f – 2D co-ordinates
type Vec2f struct {
	x float32
	y float32
}

type colourPolymorhicPrey struct {
	pos       Vec2f
	movS      float32
	movA      float32
	heading   float32
	lifetime  int32
	hunger    int32
	fertility int32 //	interval measurement between birth and sex
	gravid    bool  //	i.e. pregnant
	colour    ColRGB
}

type visualPredator struct {
	pos       Vec2f
	movS      float32
	movA      float32
	heading   float32
	lifetime  int32
	hunger    int32
	fertility int32 //	interval measurement between birth and sex
	gravid    bool  //	i.e. pregnant
	visRange  float32
	visAcuity float32
	imprint   ColRGB
}

// AgentActions interface for general agent behaviours
type AgentActions interface {
	Rotate(axis Vec2f, angle float32)
	Move()
}
