/************************
Auther: Tai Won Chung
Description: All logic for minesweeper board
************************/

package board

import (
	"container/list"
	"math/rand"
	"strconv"
	"time"
)

//region struct declarations
type Board struct {
	GameBoard    [][]int
	__innerBoard [][]int
	Width        int
	Height       int
	Nmines       int
	GameState    int
}

type pCell struct {
	row  int
	col  int
	stat int
}

//endregion struct declaration

//region const declarations
// consts for cell status for __innerBoard

const (
	EMPTY  = iota //a cell is empty if its value is 0
	NUMBER = iota //a cell is a number if its value is >=1
	MINE   = -1   //a cell has a mine if its value is -1
)

// consts for cell status for GameBoard (user-facing)
const (
	COVERED          = iota
	COVERED_FLAGGED  = iota
	UNCOVERED_EMPTY  = iota
	UNCOVERED_NUMBER = iota
	UNCOVERED_MINE   = iota
)

// consts for game state
const (
	ALIVE   = iota
	VICTORY = iota
	DEFEAT  = iota
)

//endregion const declarations

//region Board related functions
func NewBoard(width int, height int, numberOfMines int) *Board {
	b := new(Board)
	b.Width = width
	b.Height = height
	b.Nmines = numberOfMines
	b.GameState = ALIVE
	b.GameBoard = [][]int{}
	for i := 0; i < height; i++ {
		row := []int{}
		for j := 0; j < width; j++ {
			row = append(row, COVERED)
		}
		b.GameBoard = append(b.GameBoard, row)
	}
	b.__innerBoard = [][]int{}
	for i := 0; i < height; i++ {
		row := []int{}
		for j := 0; j < width; j++ {
			row = append(row, EMPTY)
		}
		b.__innerBoard = append(b.__innerBoard, row)
	}
	b.GenerateMines(width, height, numberOfMines)
	return b
}

func (b *Board) GenerateMines(width int, height int, numberOfMines int) {
	rand.Seed(time.Now().UnixNano())
	cnt := 0
	for cnt < numberOfMines {
		row := rand.Intn(height)
		col := rand.Intn(width)
		if b.__innerBoard[row][col] != MINE {
			b.__innerBoard[row][col] = MINE
			for i := row - 1; i <= row+1; i++ {
				for j := col - 1; j <= col+1; j++ {
					if i >= 0 && j >= 0 && i < height && j < width {
						if b.__innerBoard[i][j] != MINE {
							b.__innerBoard[i][j]++
						}
					}
				}
			}
			cnt++
		}
	}
}

func (b *Board) UncoverCell(row int, col int) {
	queue := list.New()
	cell := newUncoveredCell(row, col, b)
	b.GameBoard[row][col] = cell.stat
	defer b.CheckAndUpdateGameState()
	if cell.stat == UNCOVERED_MINE {
		b.GameState = DEFEAT
		return
	}
	queue.PushBack(cell)
	for elm := queue.Front(); elm != nil; elm = queue.Front() {
		cell = elm.Value.(*pCell)
		b.GameBoard[cell.row][cell.col] = cell.stat
		if cell.stat == UNCOVERED_EMPTY {
			for i := cell.row - 1; i <= cell.row+1; i++ {
				for j := cell.col - 1; j <= cell.col+1; j++ {
					if i >= 0 && j >= 0 && i < b.Height && j < b.Width {
						if tempCell := newUncoveredCell(i, j, b); tempCell != nil {
							queue.PushBack(tempCell)
						}
					}
				}
			}
		}
		_ = queue.Remove(elm)
	}
}

func (b *Board) CheckAndUpdateGameState() {
	if b.GameState == DEFEAT {
		return
	}
	cnt := 0
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			switch b.GameBoard[i][j] {
			case COVERED, COVERED_FLAGGED:
				cnt++
			}
		}
	}
	if cnt == b.Nmines {
		b.GameState = VICTORY
	}
}

func (b *Board) ToggleCellFlag(row int, col int) {
	if b.GameBoard[row][col] == COVERED {
		b.GameBoard[row][col] = COVERED_FLAGGED
	} else if b.GameBoard[row][col] == COVERED_FLAGGED {
		b.GameBoard[row][col] = COVERED
	}
}

func (b *Board) RenderBoard() [][]rune {
	renderedBoard := [][]rune{}
	for i := 0; i < b.Height; i++ {
		row := []rune{}
		for j := 0; j < b.Width; j++ {
			switch b.GameBoard[i][j] {
			case COVERED:
				row = append(row, '-')
			case COVERED_FLAGGED:
				row = append(row, '+')
			case UNCOVERED_EMPTY:
				row = append(row, '.')
			case UNCOVERED_NUMBER:
				num := []rune(strconv.FormatInt(int64(b.__innerBoard[i][j]), 10))
				row = append(row, num[0])
			case UNCOVERED_MINE:
				row = append(row, '*')
			}
		}
		renderedBoard = append(renderedBoard, row)
	}
	return renderedBoard
}

//endregion Board related functions

//region Cell related functions
func newCell(row int, col int, stat int) *pCell {
	cell := new(pCell)
	cell.row = row
	cell.col = col
	cell.stat = stat
	return cell
}

func newUncoveredCell(row int, col int, board *Board) *pCell {
	if board.GameBoard[row][col] == COVERED {
		switch {
		case board.__innerBoard[row][col] == MINE:
			return newCell(row, col, UNCOVERED_MINE)
		case board.__innerBoard[row][col] == EMPTY:
			return newCell(row, col, UNCOVERED_EMPTY)
		default:
			return newCell(row, col, UNCOVERED_NUMBER)
		}
	} else {
		return nil
	}
}

//endregion Cell related functions
