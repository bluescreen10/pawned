package main

import (
	"context"
	"math"
)

type evaluation struct {
	depth int
	best  string
	score int
}

func SearchBestMove(ctx context.Context, p *Position) chan evaluation {
	ch := make(chan evaluation)

	go func() {
		depth := 5
		eval, m := minmax(p, math.MinInt, math.MaxInt, depth)
		ch <- evaluation{
			depth: depth,
			best:  m.String(),
			score: eval,
		}
		close(ch)
	}()
	return ch
}

func minmax(p *Position, alpha, beta, depth int) (int, Move) {
	if depth == 0 {
		return eval(p), Move(0)
	}
	moves := make([]Move, 0, 218)
	moves, inCheck := LegalMoves(moves, p)

	if inCheck && len(moves) == 0 {
		if p.Active() == White {
			return math.MinInt, Move(0)
		} else {
			return math.MaxInt, Move(0)
		}
	}

	if len(moves) == 0 {
		return 0, Move(0)
	}

	if p.Active() == White {
		max := math.MinInt
		best := moves[0]
		for _, m := range moves {
			newP := *p
			Do(&newP, m)
			eval, _ := minmax(&newP, alpha, beta, depth-1)
			if eval > max {
				best = m
				max = eval
			}

			if eval >= beta {
				best = m
				break
			}

			alpha = fmax(alpha, eval)
		}
		return max, best
	} else {
		min := math.MaxInt
		best := moves[0]
		for _, m := range moves {
			newP := *p
			Do(&newP, m)
			eval, _ := minmax(&newP, alpha, beta, depth-1)
			if eval <= min {
				best = m
				min = eval
			}

			if eval <= alpha {
				best = m
				break
			}

			beta = fmin(beta, eval)

		}
		return min, best
	}
}

func eval(p *Position) int {
	moves := make([]Move, 0)
	moves, inCheck := LegalMoves(moves, p)

	if inCheck && len(moves) == 0 {
		if p.Active() == White {
			return math.MinInt
		} else {
			return math.MaxInt
		}
	}

	if len(moves) == 0 {
		return 0
	}

	wPawns := p.Pieces[White][Pawn].OnesCount()
	wKnight := p.Pieces[White][Knight].OnesCount()
	wBishop := p.Pieces[White][Bishop].OnesCount()
	wRook := p.Pieces[White][Rook].OnesCount()
	wQueen := p.Pieces[White][Queen].OnesCount()

	bPawns := p.Pieces[Black][Pawn].OnesCount()
	bKnight := p.Pieces[Black][Knight].OnesCount()
	bBishop := p.Pieces[Black][Bishop].OnesCount()
	bRook := p.Pieces[Black][Rook].OnesCount()
	bQueen := p.Pieces[Black][Queen].OnesCount()

	return (wPawns + wKnight*3 + wBishop*3 + wRook*5 + wQueen*9 -
		bPawns - bKnight*3 - bBishop*3 - bRook*5 - bQueen*9) * 100
}

func fmax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func fmin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
