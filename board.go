/************************
Auther: Tai Won Chung
Description: All logic for minesweeper board
************************/

package main

import (
	"math/rand"
	"time"
)

type Board struct {
	__board [][]int
	__mines [][]bool
	width   int
	height  int
	nmines  int
}

// consts for cell status

const (
	UNKNOWN = iota
	EMPTY   = iota
	NUMBER  = iota
	MINE    = iota
)

func NewBoard(width int, height int, numberOfMines int) *Board {
	board := new(Board)
	board.width = width
	board.height = height
	board.nmines = numberOfMines
	board.__board = [][]int{}
	for i := 0; i < height; i++ {
		row := []int{}
		for j := 0; j < width; j++ {
			row = append(row, UNKNOWN)
		}
		board.__board = append(board.__board, row)
	}
	board.__mines = [][]bool{}
	for i := 0; i < height; i++ {
		row := []bool{}
		for j := 0; j < width; j++ {
			row = append(row, false)
		}
		board.__mines = append(board.__mines, row)
	}
	GenerateMines(width, height, numberOfMines, board)
	return board
}

func GenerateMines(width int, height int, numberOfMines int, board *Board) {
	rand.Seed(time.Now().UnixNano())
	cnt := 0
	for cnt < numberOfMines {
		i := rand.Intn(height)
		j := rand.Intn(width)
		if !board.__mines[i][j] {
			board.__mines[i][j] = true
			cnt++
		}
	}
}
