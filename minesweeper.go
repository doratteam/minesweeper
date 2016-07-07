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
	fmt.Println(b.__board)
	fmt.Println(b.__mines)
}
