package solver

import (
	"github.com/doratteam/minesweeper/board"
	"testing"
)

func TestGameStateInit(t *testing.T) {
	b := board.NewBoard(9, 9, 10)
	gs := NewGameState(b)
	if gs.GraphBoard == nil {
		t.Error("Proper GameState not initialized")
	}
	for i := 1; i < 8; i++ {
		for j := 1; j < 8; j++ {
			cell := gs.CellMap[i][j]
			connCell := gs.GraphBoard[cell]
			if connCell.Cardinality() != 8 {
				t.Error("Proper GameState not initialized", connCell.Cardinality())
			}
		}
	}
}
