package solver

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/doratteam/minesweeper/board"
	//"strconv"
)

type Cell struct {
	Row        int
	Col        int
	CntMines   int
	PMine      float64
	IsFlagged  bool
	IsResolved bool
	IsCovered  bool
}

type GameState struct {
	GameBoard      *board.Board
	GraphBoard     map[*Cell]mapset.Set //utility structure to help rules
	CellMap        [][]*Cell            //represents all of current assertions
	CellsToFlag    chan *Cell
	CellsToUncover chan *Cell
}

//region GameState related functions
func NewGameState(board *board.Board) *GameState {
	gs := new(GameState)
	gs.GameBoard = board
	gs.GraphBoard = make(map[*Cell]mapset.Set)
	gs.CellMap = [][]*Cell{}
	for i := 0; i < board.Height; i++ {
		cmRow := []*Cell{}
		for j := 0; j < board.Width; j++ {
			cmRow = append(cmRow, newCell(i, j, float64(board.Nmines/(board.Height*board.Width))))
		}
		gs.CellMap = append(gs.CellMap, cmRow)
	}
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			initCellTruths(gs, i, j)
		}
	}
	gs.CellsToFlag = make(chan *Cell, 50)
	gs.CellsToUncover = make(chan *Cell, 50)
	return gs
}

func (gs *GameState) StartCellServers() {
	go cellFlaggingServer(gs)
	go cellUncoveringServer(gs)
}

func (gs *GameState) EndCellServers() {
	close(gs.CellsToFlag)
	close(gs.CellsToUncover)
}

func cellFlaggingServer(gs *GameState) {
	for cell := range gs.CellsToFlag {
		go gs.FlagCell(cell)
	}
}

func cellUncoveringServer(gs *GameState) {
	for cell := range gs.CellsToUncover {
		go gs.UncoverCell(cell)
	}
}

func initCellTruths(gs *GameState, row int, col int) {
	cell := gs.CellMap[row][col]
	connCells := mapset.NewSet()
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if !isOB(i, j, gs.GameBoard.Height, gs.GameBoard.Width) && !(i == row && j == col) {
				connCells.Add(gs.CellMap[i][j])
			}
		}
	}
	gs.GraphBoard[cell] = connCells
}

func isOB(row int, col int, height int, width int) bool {
	return row < 0 || col < 0 || row >= height || col >= width
}

func (gs *GameState) MarkCellResolved(cell *Cell) {
	row := cell.Row
	col := cell.Col
	cell.IsResolved = true
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if !isOB(i, j, gs.GameBoard.Height, gs.GameBoard.Width) && !(i == row && j == col) {
				adjCell := gs.CellMap[i][j]
				gs.GraphBoard[adjCell].Remove(cell)
			}
		}
	}
	delete(gs.GraphBoard, cell)
}

func (gs *GameState) FlagCell(cell *Cell) {
	cell.IsFlagged = true
	gs.GameBoard.FlagCell(cell.Row, cell.Col)
}

func (gs *GameState) UncoverCell(cell *Cell) {
	cell.IsCovered = false
	gs.GameBoard.UncoverCell(cell.Row, cell.Col)
}

//endregion

//region Cell related functions
func newCell(row int, col int, pmine float64) *Cell {
	c := new(Cell)
	c.Row = row
	c.Col = col
	c.CntMines = 0
	c.PMine = pmine
	c.IsFlagged = false
	c.IsResolved = false
	c.IsCovered = true
	return c
}

//endregion
