package model

import (
	"errors"
	"math"
)

// CPPbehaviour – set of actions ColourPolymorhicPrey agents will perform.
type CPPbehaviour interface {
	Turn(float64)
	Move()
	//Spawn()
	//Death() to be implemented
}

// Turn updates dir𝚯 and dir vector to the new heading offset by 𝚯
func (cpp *ColourPolymorhicPrey) Turn(𝚯 float64) {
	newHeading := UnitAngle(cpp.dir𝚯 + 𝚯)
	cpp.dir[x] = math.Cos(newHeading)
	cpp.dir[y] = math.Sin(newHeading)
	cpp.dir𝚯 = newHeading
}

// Move updates the agent's position if it doesn't encounter any errors.
func (cpp *ColourPolymorhicPrey) Move() error {
	var posOffset, newPos Vector
	var err error
	posOffset, err = VecScalarMultiply(cpp.dir, cpp.vsr)
	if err != nil {
		return errors.New("agent move failed: " + err)
	}
	newPos, err = VecAddition(cpp.pos, posOffset)
	if err != nil {
		return errors.New("agent move failed: " + err)
	}
	cpp.pos = newPos
}
