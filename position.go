package main

import (
	"fmt"
	"strconv"
	"strings"
)

// var statePool = sync.Pool{
// 	New: func() any {
// 		return &State{}
// 	},
// }

// type State struct {
// 	WhiteToMove bool

// 	// Castling rights
// 	WhiteQueenSideCastle bool
// 	WhiteKingSideCastle  bool
// 	BlackQueenSideCastle bool
// 	BlackKingSideCastle  bool

// 	// En passant file
// 	EnPassantFile uint8

// 	// Moves
// 	HalfMoves uint8
// 	FullMoves uint16

// 	// Previous state
// 	LastMove *Move
// 	Capture  Piece
// 	Previous *State
// }

// func (s *State) Reset() {
// 	s.WhiteToMove = false
// 	s.WhiteKingSideCastle = false
// 	s.WhiteQueenSideCastle = false
// 	s.BlackKingSideCastle = false
// 	s.BlackQueenSideCastle = false
// 	s.EnPassantFile = 0
// 	s.HalfMoves = 0
// 	s.FullMoves = 0
// 	s.Previous = nil
// }

type Color uint8

const (
	White Color = iota
	Black
)

type Piece uint8

const (
	Pawn Piece = iota
	Knight
	Bishop
	Rook
	Queen
	King
	Empty
)

type CastlingRights uint8

const (
	WhiteKingSideCastle CastlingRights = 1 << iota
	WhiteQueenSideCastle
	BlackKingSideCastle
	BlackQueenSideCastle
)

type Position struct {
	Pieces    [Color(2)][Piece(6)]BitBoard
	AllPieces [Color(2)]BitBoard
	Occupied  BitBoard
	//State     *State

	WhiteToMove     bool
	EnPassantSquare Square
	CastlingRights  CastlingRights
	HalfMoves       uint8
	FullMoves       uint16
}

func Parse(fen string) (Position, error) {
	var pos Position

	parts := strings.Split(strings.TrimSpace(fen), " ")
	if len(parts) < 6 {
		return pos, fmt.Errorf("invalid fen: %s", fen)
	}

	bit := BitBoard(1)

	for _, row := range strings.Split(parts[0], "/") {
		for _, char := range row {
			switch char {
			case 'P':
				pos.Pieces[White][Pawn] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'N':
				pos.Pieces[White][Knight] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'B':
				pos.Pieces[White][Bishop] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'R':
				pos.Pieces[White][Rook] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'Q':
				pos.Pieces[White][Queen] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'K':
				pos.Pieces[White][King] |= bit
				pos.AllPieces[White] |= bit
				pos.Occupied |= bit
			case 'p':
				pos.Pieces[Black][Pawn] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case 'n':
				pos.Pieces[Black][Knight] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case 'b':
				pos.Pieces[Black][Bishop] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case 'r':
				pos.Pieces[Black][Rook] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case 'q':
				pos.Pieces[Black][Queen] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case 'k':
				pos.Pieces[Black][King] |= bit
				pos.AllPieces[Black] |= bit
				pos.Occupied |= bit
			case '1', '2', '3', '4', '5', '6', '7', '8':
				bit <<= uint(char - '1')
			default:
				return Position{}, fmt.Errorf("invalid piece: %c", char)
			}
			bit <<= 1
		}
	}

	// pos.State = statePool.Get().(*State)
	// pos.State.Reset()

	if parts[1] == "w" || parts[1] == "W" {
		pos.WhiteToMove = true
	}

	for _, c := range parts[2] {
		switch c {
		case 'K':
			pos.CastlingRights |= WhiteKingSideCastle
		case 'Q':
			pos.CastlingRights |= WhiteQueenSideCastle
		case 'k':
			pos.CastlingRights |= BlackKingSideCastle
		case 'q':
			pos.CastlingRights |= BlackQueenSideCastle
		}
	}

	pos.EnPassantSquare = SquareFromString(parts[3])

	halfMoves, err := strconv.Atoi(parts[4])
	if err != nil {
		return pos, fmt.Errorf("invalid half moves: %s", fen)
	}
	pos.HalfMoves = uint8(halfMoves)

	fullMoves, err := strconv.Atoi(parts[5])
	if err != nil {
		return pos, fmt.Errorf("invalid full moves: %s", fen)
	}
	pos.FullMoves = uint16(fullMoves)

	return pos, nil
}

func (p Position) Fen() string {
	fen := strings.Builder{}

	for bit := BitBoard(1); bit != 0; bit <<= 1 {
		if bit&File_A != 0 && bit > 1 {
			fen.WriteByte('/')
		}

		if p.Pieces[White][Pawn]&bit != 0 {
			fen.WriteByte('P')
		} else if p.Pieces[Black][Pawn]&bit != 0 {
			fen.WriteByte('p')
		} else if p.Pieces[White][Knight]&bit != 0 {
			fen.WriteByte('N')
		} else if p.Pieces[Black][Knight]&bit != 0 {
			fen.WriteByte('n')
		} else if p.Pieces[White][Bishop]&bit != 0 {
			fen.WriteByte('B')
		} else if p.Pieces[Black][Bishop]&bit != 0 {
			fen.WriteByte('b')
		} else if p.Pieces[White][Rook]&bit != 0 {
			fen.WriteByte('R')
		} else if p.Pieces[Black][Rook]&bit != 0 {
			fen.WriteByte('r')
		} else if p.Pieces[White][Queen]&bit != 0 {
			fen.WriteByte('Q')
		} else if p.Pieces[Black][Queen]&bit != 0 {
			fen.WriteByte('q')
		} else if p.Pieces[White][King]&bit != 0 {
			fen.WriteByte('K')
		} else if p.Pieces[Black][King]&bit != 0 {
			fen.WriteByte('k')
		} else {
			empty := 1
			bit <<= 1
			for ; bit != 0; bit, empty = bit<<1, empty+1 {
				if p.Occupied&bit != 0 || bit&File_A != 0 {
					bit >>= 1
					break
				}
			}

			fen.WriteByte('0' + byte(empty))
		}
	}

	fen.WriteByte(' ')

	stm, _ := p.SideToMove()
	if stm == White {
		fen.WriteByte('w')
	} else {
		fen.WriteByte('b')
	}

	fen.WriteByte(' ')

	if p.CanWhiteCastleKingSide() || p.CanWhiteCastleQueenSide() || p.CanBlackCastleKingSide() || p.CanBlackCastleQueenSide() {
		if p.CanWhiteCastleKingSide() {
			fen.WriteByte('K')
		}

		if p.CanWhiteCastleQueenSide() {
			fen.WriteByte('Q')
		}

		if p.CanBlackCastleKingSide() {
			fen.WriteByte('k')
		}

		if p.CanBlackCastleQueenSide() {
			fen.WriteByte('q')
		}
	} else {
		fen.WriteString("-")
	}

	fen.WriteByte(' ')
	fen.WriteString(p.EnPassantSquare.String())
	fen.WriteString(fmt.Sprintf(" %d %d", p.HalfMoves, p.FullMoves))

	return fen.String()
}

func (p Position) String() string {
	builder := strings.Builder{}

	builder.WriteString("+---+---+---+---+---+---+---+---+\n")
	bit := BitBoard(1)
	for rank := 7; rank >= 0; rank-- {
		builder.WriteByte('|')
		for file := 0; file < 8; file++ {
			if p.Pieces[White][Pawn]&bit != 0 {
				builder.WriteString(" P |")
			} else if p.Pieces[Black][Pawn]&bit != 0 {
				builder.WriteString(" p |")
			} else if p.Pieces[White][Knight]&bit != 0 {
				builder.WriteString(" N |")
			} else if p.Pieces[Black][Knight]&bit != 0 {
				builder.WriteString(" n |")
			} else if p.Pieces[White][Bishop]&bit != 0 {
				builder.WriteString(" B |")
			} else if p.Pieces[Black][Bishop]&bit != 0 {
				builder.WriteString(" b |")
			} else if p.Pieces[White][Rook]&bit != 0 {
				builder.WriteString(" R |")
			} else if p.Pieces[Black][Rook]&bit != 0 {
				builder.WriteString(" r |")
			} else if p.Pieces[White][Queen]&bit != 0 {
				builder.WriteString(" Q |")
			} else if p.Pieces[Black][Queen]&bit != 0 {
				builder.WriteString(" q |")
			} else if p.Pieces[White][King]&bit != 0 {
				builder.WriteString(" K |")
			} else if p.Pieces[Black][King]&bit != 0 {
				builder.WriteString(" k |")
			} else {
				builder.WriteString("   |")
			}
			bit <<= 1
		}
		builder.WriteString(fmt.Sprintf(" %d\n", rank+1))
		builder.WriteString("+---+---+---+---+---+---+---+---+\n")
	}
	builder.WriteString("  a   b   c   d   e   f   g   h\n")
	return builder.String()
}

func (p Position) SideToMove() (Color, Color) {
	if p.WhiteToMove {
		return White, Black
	}
	return Black, White
}

func (p Position) IsEnPassant() bool {
	return p.EnPassantSquare != SQ_NULL
}

func (p Position) EnPassantFile() BitBoard {
	_, file := p.EnPassantSquare.RankAndFile()
	switch file {
	case 0:
		return File_A
	case 1:
		return File_B
	case 2:
		return File_C
	case 3:
		return File_D
	case 4:
		return File_E
	case 5:
		return File_F
	case 6:
		return File_G
	case 7:
		return File_H
	default:
		return BitBoard(0)
	}
}

func (p *Position) CanWhiteCastleKingSide() bool {
	return p.CastlingRights&WhiteKingSideCastle != 0
}

func (p *Position) CanWhiteCastleQueenSide() bool {
	return p.CastlingRights&WhiteQueenSideCastle != 0
}

func (p *Position) CanBlackCastleKingSide() bool {
	return p.CastlingRights&BlackKingSideCastle != 0
}

func (p *Position) CanBlackCastleQueenSide() bool {
	return p.CastlingRights&BlackQueenSideCastle != 0
}

func (p Position) Do(m Move) Position {

	new := p

	us, them := new.SideToMove()
	new.removeAll(them, m.To)
	new.remove(m.Piece, us, m.From)

	switch m.Type {
	case EnPassant:
		if us == White {
			new.remove(Pawn, them, m.To+8)
		} else {
			new.remove(Pawn, them, m.To-8)
		}
		new.put(Pawn, us, m.To)
	case PromotionToKnight:
		new.put(Knight, us, m.To)
	case PromotionToBishop:
		new.put(Bishop, us, m.To)
	case PromotionToRook:
		new.put(Rook, us, m.To)
	case PromotionToQueen:
		new.put(Queen, us, m.To)
	case Castle:
		switch m.To {
		case SQ_G1:
			new.put(King, White, SQ_G1)
			new.put(Rook, White, SQ_F1)
			new.remove(Rook, White, SQ_H1)
		case SQ_C1:
			new.put(King, White, SQ_C1)
			new.put(Rook, White, SQ_D1)
			new.remove(Rook, White, SQ_A1)
		case SQ_G8:
			new.put(King, Black, SQ_G8)
			new.put(Rook, Black, SQ_F8)
			new.remove(Rook, Black, SQ_H8)
		case SQ_C8:
			new.put(King, Black, SQ_C8)
			new.put(Rook, Black, SQ_D8)
			new.remove(Rook, Black, SQ_A8)
		}
	default:
		new.put(m.Piece, us, m.To)
	}

	new.EnPassantSquare = SQ_NULL

	if m.Piece == Pawn {
		if m.From-m.To == 16 {
			new.EnPassantSquare = m.To + 8
		}

		if m.To-m.From == 16 {
			new.EnPassantSquare = m.To - 8
		}
	}

	if new.CastlingRights != 0 {
		if m.Type == Castle && us == White {
			new.CastlingRights &^= WhiteKingSideCastle
			new.CastlingRights &^= WhiteQueenSideCastle
		}

		if m.Type == Castle && us == Black {
			new.CastlingRights &^= BlackKingSideCastle
			new.CastlingRights &^= BlackQueenSideCastle
		}

		if m.From == SQ_A1 || m.To == SQ_A1 {
			new.CastlingRights &^= WhiteQueenSideCastle
		}

		if m.From == SQ_H1 || m.To == SQ_H1 {
			new.CastlingRights &^= WhiteKingSideCastle
		}

		if m.From == SQ_A8 || m.To == SQ_A8 {
			new.CastlingRights &^= BlackQueenSideCastle
		}

		if m.From == SQ_H8 || m.To == SQ_H8 {
			new.CastlingRights &^= BlackKingSideCastle
		}
	}

	new.WhiteToMove = !p.WhiteToMove
	new.HalfMoves++

	if us == Black {
		new.FullMoves++
	}

	return new
}

func (p *Position) put(piece Piece, color Color, sq Square) {
	bit := BitBoard(1) << sq
	p.Occupied |= bit
	p.AllPieces[color] |= bit
	p.Pieces[color][piece] |= bit
}

func (p *Position) get(color Color, sq Square) Piece {
	bit := BitBoard(1) << sq
	if p.Pieces[color][Pawn]&bit != 0 {
		return Pawn
	} else if p.Pieces[color][Knight]&bit != 0 {
		return Knight
	} else if p.Pieces[color][Bishop]&bit != 0 {
		return Bishop
	} else if p.Pieces[color][Rook]&bit != 0 {
		return Rook
	} else if p.Pieces[color][Queen]&bit != 0 {
		return Queen
	} else if p.Pieces[color][King]&bit != 0 {
		return King
	} else {
		return Empty
	}
}

func (p *Position) remove(piece Piece, color Color, sq Square) {
	bit := BitBoard(1) << sq
	p.Occupied &^= bit
	p.AllPieces[color] &^= bit
	p.Pieces[color][piece] &^= bit
}

func (p *Position) removeAll(color Color, sq Square) {
	bit := BitBoard(1) << sq
	p.Occupied &^= bit
	p.AllPieces[color] &^= bit
	p.Pieces[color][Pawn] &^= bit
	p.Pieces[color][Knight] &^= bit
	p.Pieces[color][Bishop] &^= bit
	p.Pieces[color][Rook] &^= bit
	p.Pieces[color][Queen] &^= bit
	//p.Pieces[color][King] &^= bit
}
