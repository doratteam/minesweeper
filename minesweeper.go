/*********************************
Auther: Tai Won Chung
Description: This is the main driver for minesweeper game
*********************************/

package main

import (
	"fmt"
)

func main() {
	b := NewBoard(9, 9, 10)
	fmt.Printf("created board of %dx%d\n", b.width, b.height)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%2d ", b.gameBoard[i][j])
		}
		fmt.Printf("\n")
	}
	b.UncoverCell(0, 0)
	fmt.Printf("\n")
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%2d ", b.gameBoard[i][j])
		}
		fmt.Printf("\n")
	}
}
