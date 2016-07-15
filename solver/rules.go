package solver

import (
	mapset "github.com/deckarep/golang-set"
)

type Rule interface {
	Condition(*GameState) bool
	ActionToTake(*GameState)
}

type BasicUnicellRule struct {
	A    *Cell
	AdjA mapset.Set
}

func (r BasicUnicellRule) Condition(gs *GameState) bool {
	nmines := r.A.CntMines
	ncovered := 0
	for c := range r.AdjA.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered || cell.IsFlagged {
			ncovered++
		}
	}
	return nmines == ncovered
}

func (r BasicUnicellRule) ActionToTake(gs *GameState) {
	for c := range r.AdjA.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered {
			gs.CellsToFlag <- cell
		}
	}
}

type CheckResolvedRule struct {
	A    *Cell
	AdjA mapset.Set
}

func (r CheckResolvedRule) Condition(gs *GameState) bool {
	nmines := r.A.CntMines
	nflagged := 0
	for c := range r.AdjA.Iter() {
		cell := c.(*Cell)
		if cell.IsFlagged {
			nflagged++
		}
	}
	return nmines == nflagged
}

func (r CheckResolvedRule) ActionToTake(gs *GameState) {
	gs.MarkCellResolved(r.A)
	for c := range r.AdjA.Iter() {
		cell := c.(*Cell)
		if cell.IsCovered {
			gs.CellsToUncover <- cell
		}
	}
}
