package types

import (
	qtypes "github.com/rossus/quadria/common/types"
	"html/template"
)

// UITile is a minimal tile representation for the animation frames.
type UITile struct {
	Value int    `json:"value"`
	Color string `json:"color"`
}

// AnimateData is passed to the animation template.
type AnimateData struct {
	Steps   template.JS
	Turn    int
	Players []qtypes.Player
	Active  qtypes.Player
	GameURL string
}
