/************************
Auther: Tai Won Chung
Description: All logic for minesweeper board
************************/

package main

import (
	"math/rand"
	"time"
)

//region struct declarations
type Board struct {
	gameBoard    [][]int
	__innerBoard [][]int
	width        int
	height       int
	nmines       int
}

type Cell struct {
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

// consts for cell status for gameBoard (user-facing)
const (
	COVERED          = iota
	COVERED_FLAGGED  = iota
	UNCOVERED_EMPTY  = iota
	UNCOVERED_NUMBER = iota
	UNCOVERED_MINE   = iota
)

//endregion const declarations

//region Board related functions
func NewBoard(width int, height int, numberOfMines int) *Board {
	b := new(Board)
	b.width = width
	b.height = height
	b.nmines = numberOfMines
	b.gameBoard = [][]int{}
	for i := 0; i < height; i++ {
		row := []int{}
		for j := 0; j < width; j++ {
			row = append(row, COVERED)
		}
		b.gameBoard = append(b.gameBoard, row)
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
	queue := []*Cell{}
	queue = append(queue, NewUncoveredCell(row, col, b))
}

//endregion Board related functions

//region Cell related functions
func NewCell(row int, col int, stat int) *Cell {
	cell := new(Cell)
	cell.row = row
	cell.col = col
	cell.stat = stat
	return cell
}

func NewUncoveredCell(row int, col int, board *Board) *Cell {
	switch {
	case board.__innerBoard[row][col] == MINE:
		return NewCell(row, col, UNCOVERED_MINE)
	case board.__innerBoard[row][col] == EMPTY:
		return NewCell(row, col, UNCOVERED_EMPTY)
	default:
		return NewCell(row, col, UNCOVERED_NUMBER)
	}
}

//endregion Cell related functions
