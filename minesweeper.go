/*********************************
Auther: Tai Won Chung
Description: This is the main driver for minesweeper game
*********************************/

package main

import (
	"fmt"
	"github.com/doratteam/minesweeper/board"
)

func main() {
	b := board.NewBoard(9, 9, 10)
	fmt.Printf("created board of %dx%d\n", b.Width, b.Height)
	rb := b.RenderBoard()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%c ", rb[i][j])
		}
		fmt.Printf("\n")
	}
	b.UncoverCell(0, 0)
	fmt.Printf("\n")
	rb = b.RenderBoard()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%c ", rb[i][j])
		}
		fmt.Printf("\n")
	}
}
