package types

import qtypes "github.com/rossus/quadria/common/types"

// GamePageData holds information passed to the game template.
type GamePageData struct {
	Tiles   [][]qtypes.Tile
	Turn    int
	Players []qtypes.Player
	Active  qtypes.Player
}
