package solver

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/doratteam/minesweeper/board"
	"strconv"
)

type Cell struct {
	Row       int
	Col       int
	CntMines  int
	IsFlagged bool
}

type GameState struct {
	GameBoard   *board.Board
	CntFlagged  int
	MappedBoard map[*Cell][]*Cell
	NumCells    mapset.Set
	EdgeCells   mapset.Set
	CellMap     [][]*Cell
}

//region GameState related functions
func NewGameState(board *board.Board) *GameState {
	gs := new(GameState)
	gs.NumCells = mapset.NewSet()
	gs.EdgeCells = mapset.NewSet()
	gs.AllCells = mapset.NewSet()
	gs.CellMap = [][]*Cell{}
	for i := 0; i < board.Height; i++ {
		cmRow = []*Cell{}
		for j := 0; j < board.Width; j++ {
			c := newCell(i, j)
			gs.AllCells.Add(c)
			cmRow = append(cmRow, c)
		}
		gs.CellMap = append(gs.CellMap, cmRow)
	}
	return gs
}

func (g *GameState) UpdateState() {
	rb := g.GameBoard.RenderBoard()
	for i := 0; i < g.GameBoard.Height; i++ {
		for j := 0; j < g.GameBoard.Width; j++ {
			switch {
			case rb[i][j] == '-':
				if gs.CellMap[i][j].IsFlagged {
					gs.CellMap[i][j].IsFlagged = false
				}
			case rb[i][j] == '+':
				if !gs.CellMap[i][j].IsFlagged {
					gs.CellMap[i][j].IsFlagged = true
				}
			default:
				break
			}
		}
	}
}

//endregion

//region Cell related functions
func newCell(row int, col int) *Cell {
	c := new(Cell)
	c.Row = row
	c.Col = col
	c.CntMines = 0
	c.IsFlagged = false
	return c
}

//endregion
