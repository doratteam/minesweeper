package solver

import (
	"github.com/doratteam/minesweeper/board"
	"testing"
)

func TestGameStateInit(t *testing.T) {
	b := board.NewBoard(9, 9, 10)
	g := NewGameState(b)
	if g.AllCells.Cardinality() != 81 {
		t.Error("Proper GameState not initialized")
	}
}
