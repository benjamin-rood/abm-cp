package render

import (
	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
)

// pos2D is the in-between translation for vector position.
type pos2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AgentRender holds the minimum required data for the data visualisation of an individual agent
type AgentRender struct {
	Pos     pos2D      `json:"position"`
	Heading float64    `json:"heading"`
	Col     colour.RGB `json:"colour"`
}

// Msg contains the draw instructions for front-end JS gfx API
type Msg struct {
	Type string        `json:"type"`
	CPP  []AgentRender `json:"render-cpp-array"`
	VP   []AgentRender `json:"render-vp-array"`
}

// Viewport holds the resolution / scale of the front-end JS gfx
type Viewport struct {
	Type          string
	Width, Height float64
}

func TranslateToViewport(ar AgentRender, v Viewport) (out AgentRender) {
	// for now we have to assume that the range of x = (-1.0,+1.0) and y = (-1.0,+1.0)
	out.Pos.X = absToView(ar.Pos.X, 1.0, v.Width)
	out.Pos.Y = absToView(ar.Pos.Y, 1.0, v.Height)
	out.Col = ar.Col
	return
}

func absToView(p float64, d float64, n float64) (view float64) {
	view = ((p + d) / (2 * d)) * n
	view = calc.ToFixed(view, 3)
	return
}
