package solver

import (
	"container/list"
	"strconv"

	mapset "github.com/deckarep/golang-set"
	"github.com/doratteam/minesweeper/board"
)

// Cell struct
type Cell struct {
	Row        int
	Col        int
	CntMines   int
	PMine      float64
	IsFlagged  bool
	IsResolved bool
	IsCovered  bool
}

//GameState (a.k.a. Truths we know about the game)
type GameState struct {
	GameBoard      *board.Board
	GraphBoard     map[*Cell]mapset.Set //utility structure to help rules
	CellMap        [][]*Cell            //represents all of current assertions
	CellsToFlag    chan *Cell
	CellsToUncover chan *Cell
}

//region GameState related functions

//NewGameState inits a GameState given a board
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

//StartCellServers Starts cell servers
func (gs *GameState) StartCellServers() {
	go cellFlaggingServer(gs)
	go cellUncoveringServer(gs)
}

//EndCellServers Closes channels to stop cell servers
func (gs *GameState) EndCellServers() {
	close(gs.CellsToFlag)
	close(gs.CellsToUncover)
}

//cellFlaggingServer Server that flags cells
func cellFlaggingServer(gs *GameState) {
	for cell := range gs.CellsToFlag {
		go gs.FlagCell(cell)
	}
}

//cellUncoveringServer Server that uncvers cells
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

//MarkCellResolved marks a cell as resolved and removes it from GraphBoard
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

//FlagCell flags a given cell
func (gs *GameState) FlagCell(cell *Cell) {
	cell.IsFlagged = true
	gs.GameBoard.FlagCell(cell.Row, cell.Col)
}

//UncoverCell uncovers a given cell (cascading)
func (gs *GameState) UncoverCell(cell *Cell) {
	cell.IsCovered = false
	gs.GameBoard.UncoverCell(cell.Row, cell.Col)
	gs.CascadeUpdateCell(cell)
}

//CascadeUpdateCell updates truths about cells (cascading)
func (gs *GameState) CascadeUpdateCell(cell *Cell) {
	queue := list.New()
	queue.PushBack(cell)
	visited := mapset.NewSet()
	visited.Add(cell)
	for elm := queue.Front(); elm != nil; elm = queue.Front() {
		cell = elm.Value.(*Cell)
		isUpdated := gs.UpdateCell(cell)
		if isUpdated {
			for i := cell.Row - 1; i <= cell.Row+1; i++ {
				for j := cell.Col - 1; j <= cell.Col+1; j++ {
					if !isOB(i, j, gs.GameBoard.Height, gs.GameBoard.Width) && !visited.Contains(gs.CellMap[i][j]) {
						queue.PushBack(gs.CellMap[i][j])
						visited.Add(gs.CellMap[i][j])
					}
				}
			}
		}
		_ = queue.Remove(elm)
	}
}

//UpdateCell updates truths about a given cell (according to its rendering)
func (gs *GameState) UpdateCell(cell *Cell) bool {
	isUpdated := false
	rc := gs.GameBoard.RenderCell(cell.Row, cell.Col)
	switch rc {
	case '.':
		if cell.IsCovered {
			cell.IsCovered = false
			isUpdated = true
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if cell.IsCovered {
			cell.IsCovered = false
			i64, _ := strconv.ParseInt(string(rc), 10, 0)
			cell.CntMines = int(i64)
			isUpdated = true
		}
	}
	return isUpdated
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
