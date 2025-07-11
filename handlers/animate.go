package handlers

import (
	"github.com/rossus/quadria/board"
	qtypes "github.com/rossus/quadria/common/types"
	"github.com/rossus/quadria/session"

	"github.com/rossus/codex-gen-quadria-ui/types"
)

// addToNext merges checked actions with the slice of next actions.
func addToNext(next []qtypes.NextAction, checked [][3]int) []qtypes.NextAction {
	for i := 0; i < len(checked); i++ {
		for j := 0; j < len(next); j++ {
			if next[j].X == checked[i][0] && next[j].Y == checked[i][1] {
				next[j].Amount += checked[i][2]
				i++
				break
			}
		}
		if i < len(checked) {
			next = append(next, qtypes.NextAction{X: checked[i][0], Y: checked[i][1], Amount: checked[i][2]})
		}
	}
	return next
}

// checkOverlap returns actions triggered when a tile exceeds its neighbours.
func checkOverlap(s *session.Session, x, y int) [][3]int {
	var actions [][3]int
	tile := s.Board.GetTile(x, y)
	if tile.Value <= tile.Neighbours {
		return nil
	}
	grade := tile.Value / tile.Neighbours
	loose := -tile.Neighbours * grade
	actions = append(actions, [3]int{x, y, loose})
	if x != 0 {
		actions = append(actions, [3]int{x - 1, y, grade})
	}
	if x != len(s.Board.GetTiles())-1 {
		actions = append(actions, [3]int{x + 1, y, grade})
	}
	if y != 0 {
		actions = append(actions, [3]int{x, y - 1, grade})
	}
	if y != len(s.Board.GetTiles())-1 {
		actions = append(actions, [3]int{x, y + 1, grade})
	}
	return actions
}

// goSub performs one subturn and returns actions for the next subturn.
func goSub(s *session.Session, current []qtypes.NextAction) ([]qtypes.NextAction, bool) {
	next := make([]qtypes.NextAction, 0)
	for i := 0; i < len(current); i++ {
		oldVal := s.Board.GetTile(current[i].X, current[i].Y).Value
		s.Board.ChangeTileState(current[i].X, current[i].Y, oldVal+current[i].Amount)
		s.Game.ActionDone(current[i].X, current[i].Y, oldVal, oldVal+current[i].Amount)
		if s.Board.CheckDomination() {
			return nil, true
		}
		next = addToNext(next, checkOverlap(s, current[i].X, current[i].Y))
	}
	return next, false
}

// snapshot returns a simplified board representation.
func snapshot(b *board.Board) [][]types.UITile {
	tiles := b.GetTiles()
	snap := make([][]types.UITile, len(tiles))
	for y, row := range tiles {
		snap[y] = make([]types.UITile, len(row))
		for x, t := range row {
			snap[y][x] = types.UITile{Value: t.Value, Color: t.Player.Color}
		}
	}
	return snap
}

// runMove executes a move and records board states after each subturn.
func runMove(s *session.Session, x, y int) ([][][]types.UITile, bool) {
	frames := make([][][]types.UITile, 0)
	next := []qtypes.NextAction{{X: x, Y: y, Amount: 1}}
	var done bool
	for {
		next, done = goSub(s, next)
		frames = append(frames, snapshot(s.Board))
		if done {
			return frames, true
		}
		if len(next) == 0 {
			s.Game.NextTurn()
			break
		}
		s.Game.NextSubTurn()
	}
	return frames, false
}
