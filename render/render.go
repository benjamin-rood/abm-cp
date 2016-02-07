package render

import (
	"github.com/benjamin-rood/abm-cp/calc"
	"github.com/benjamin-rood/abm-cp/colour"
)

// Pos2D is the pixel 2D projection from vector position.
type Pos2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AgentRender holds the minimum required data for the data visualisation of an individual agent
type AgentRender struct {
	Type    string `json:"agent-type"`
	Pos2D   `json:"position"`
	Heading float64       `json:"heading"`
	Colour  colour.RGB256 `json:"colour"`
}

// DrawList contains the draw instructions for front-end JS gfx API
type DrawList struct {
	CPP       []AgentRender `json:"cpp"`
	VP        []AgentRender `json:"vp"`
	BG        colour.RGB256 `json:"bg"`
	CppPop    string        `json:"cpp-pop-string"` // this and the next three entries to be displayed in a little box
	VpPop     string        `json:"vp-pop-string"`
	TurnCount string        `json:"turncount-string"`
}

// Viewport holds the resolution / scale of the front-end JS gfx
type Viewport struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// TranslateToViewport takes the absolute model coordinates of the agent's position data and translates (and scales) them to the pixel coordinates of the Viewport v.
func (ar *AgentRender) TranslateToViewport(v Viewport, dw float64, dh float64) {
	ar.X = absToView(ar.X, dw, v.Width)
	ar.Y = absToView(ar.Y, dh, v.Height)
	return
}

func absToView(p float64, d float64, n float64) (view float64) {
	view = ((p + d) / (2 * d)) * n
	view = calc.ToFixed(view, 3)
	return
}
