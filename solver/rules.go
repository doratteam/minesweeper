package solver

import (
	mapset "github.com/deckarep/golang-set"
)

type Rule interface {
	Condition(*GameState) bool
	ActionToTake(*GameState)
}

type BasicUnicellRule struct {
	A     *Cell
	Adj_A mapset.Set
}

func (r BasicUnicellRule) Condition(gs *GameState) bool {
	n_mines := r.A.CntMines
	n_covered := 0
	for c := range r.Adj_A.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered {
			n_covered++
		}
	}
	return n_mines == n_covered
}

func (r BasicUnicellRule) ActionToTake(gs *GameState) {
	for c := range r.Adj_A.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered {
			gs.CellsToFlag <- cell
		}
	}
}

type CheckResolvedRule struct {
	A     *Cell
	Adj_A mapset.Set
}

func (r CheckResolvedRule) Condition(gs *GameState) bool {
	n_mines := r.A.CntMines
	n_flagged := 0
	for c := range r.Adj_A.Iter() {
		cell := c.(*Cell)
		if cell.IsFlagged {
			n_flagged++
		}
	}
	return n_mines == n_flagged
}

func (r CheckResolvedRule) ActionToTake(gs *GameState) {
	gs.MarkCellResolved(r.A)
	for c := range r.Adj_A.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered {
			gs.CellsToUncover <- cell
		}
	}
}
