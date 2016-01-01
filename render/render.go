package render

import "github.com/benjamin-rood/abm-colour-polymorphism/colour"

// pos2D is the in-between translation for vector position.
type pos2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// CppRender holds the minimum required data for the data visualisation of an individual ColourPolymorhicPrey agent.
type CppRender struct {
	Pos pos2D      `json:"position"`
	Col colour.RGB `json:"colouration"`
}

// VpRender holds the minimum required data for the data visualisation of an individual VisualPredator agent.
type VpRender struct {
	Pos     pos2D      `json:"position"`
	Heading float64    `json:"heading"`
	Col     colour.RGB `json:"imprinting"`
}

// Msg contains the draw instructions for front-end JS gfx API
type Msg struct {
	Type string      `json:"type"`
	CPP  []CppRender `json:"render-cpp-array"`
	VP   []VpRender  `json:"render-vp-array"`
}

// Viewport holds the resolution / scale of the front-end JS gfx
type Viewport struct {
}

func translateToViewport(r CppRender, v Viewport) (out CppRender) {
	return
}
