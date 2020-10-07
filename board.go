package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Move struct {
	startIndex uint8
	endIndex   uint8
}

func (m *Move) String() string {
	return fmt.Sprintf("%s to %s", indexToStringPos(m.startIndex), indexToStringPos(m.endIndex))
}

func indexFromStringLoc(pos string) (uint8, error) {
	if len(pos) != 2 {
		return 65, fmt.Errorf("position string (%s) invalid lenth.", pos)
	}
	column := pos[0] - 'a'
	row := pos[1] - '1'

	if column < 0 || column > 7 {
		return 65, fmt.Errorf("position string (%s) invalid column.", pos)
	}
	if row < 0 || row > 7 {
		return 65, fmt.Errorf("position string (%s) invalid row.", pos)
	}
	return row*8 + column, nil
}

func MoveFromStringPos(pos string) (*Move, error) {
	moves := strings.Split(pos, " to ")
	if len(moves) != 2 {
		moves = strings.Split(pos, " ")
		if len(moves) != 2 {
			return nil, fmt.Errorf("could not parse move string %s", pos)
		}
	}
	m := &Move{}
	var err error
	m.startIndex, err = indexFromStringLoc(moves[0])
	if err != nil {
		return nil, err
	}
	m.endIndex, err = indexFromStringLoc(moves[1])
	if err != nil {
		return nil, err
	}
	return m, nil
}

func indexToStringPos(index uint8) string {
	row := index / 8
	column := index % 8

	return fmt.Sprintf("%b%d", 'a'+column, row)
}

type ColorIndex int
type PieceIndex int

const (
	KingIndex   PieceIndex = 0
	QueenIndex  PieceIndex = 1
	BishopIndex PieceIndex = 2
	KnightIndex PieceIndex = 3
	RookIndex   PieceIndex = 4
	PawnIndex   PieceIndex = 5

	WhiteIndex ColorIndex = 0
	BlackIndex ColorIndex = 1
)

type Board struct {
	pieces [2][6]uint64
}

func NewStandardBoard() *Board {
	b := &Board{}
	b.pieces[WhiteIndex][KingIndex] = 1 << 4
	b.pieces[BlackIndex][KingIndex] = 1 << 60
	b.pieces[WhiteIndex][QueenIndex] = 1 << 3
	b.pieces[BlackIndex][QueenIndex] = 1 << 59
	b.pieces[WhiteIndex][BishopIndex] = 1<<2 | 1<<5
	b.pieces[BlackIndex][BishopIndex] = 1<<58 | 1<<61
	b.pieces[WhiteIndex][KnightIndex] = 1<<1 | 1<<6
	b.pieces[BlackIndex][KnightIndex] = 1<<57 | 1<<62
	b.pieces[WhiteIndex][RookIndex] = 1 | 1<<7
	b.pieces[BlackIndex][RookIndex] = 1<<56 | 1<<63

	b.pieces[WhiteIndex][PawnIndex] = 0
	for i := 8; i < 16; i++ {
		b.pieces[WhiteIndex][PawnIndex] |= 1 << i
	}
	b.pieces[BlackIndex][PawnIndex] = 0
	for i := 48; i < 56; i++ {
		b.pieces[BlackIndex][PawnIndex] |= 1 << i
	}
	return b
}

func (b *Board) GetWhitePieces() uint64 {
	return b.pieces[WhiteIndex][KingIndex] |
		b.pieces[WhiteIndex][QueenIndex] |
		b.pieces[WhiteIndex][BishopIndex] |
		b.pieces[WhiteIndex][KnightIndex] |
		b.pieces[WhiteIndex][RookIndex] |
		b.pieces[WhiteIndex][PawnIndex]
}

func (b *Board) GetBlackPieces() uint64 {
	return b.pieces[BlackIndex][KingIndex] |
		b.pieces[BlackIndex][QueenIndex] |
		b.pieces[BlackIndex][BishopIndex] |
		b.pieces[BlackIndex][KnightIndex] |
		b.pieces[BlackIndex][RookIndex] |
		b.pieces[BlackIndex][PawnIndex]
}

func (b *Board) GetAllPieces() uint64 {
	return b.GetBlackPieces() | b.GetWhitePieces()
}

func (b *Board) String() string {
	pieces := make([]byte, 64)
	for i := uint64(0); i < 64; i++ {
		pieces[i] = ' '
	}
	addPiece(b.pieces[WhiteIndex][KingIndex], 'K', pieces)
	addPiece(b.pieces[BlackIndex][KingIndex], 'k', pieces)
	addPiece(b.pieces[WhiteIndex][QueenIndex], 'Q', pieces)
	addPiece(b.pieces[BlackIndex][QueenIndex], 'q', pieces)
	addPiece(b.pieces[WhiteIndex][BishopIndex], 'B', pieces)
	addPiece(b.pieces[BlackIndex][BishopIndex], 'b', pieces)
	addPiece(b.pieces[WhiteIndex][KnightIndex], 'H', pieces)
	addPiece(b.pieces[BlackIndex][KnightIndex], 'h', pieces)
	addPiece(b.pieces[WhiteIndex][RookIndex], 'R', pieces)
	addPiece(b.pieces[BlackIndex][RookIndex], 'r', pieces)
	addPiece(b.pieces[WhiteIndex][PawnIndex], 'P', pieces)
	addPiece(b.pieces[BlackIndex][PawnIndex], 'p', pieces)

	retval := "-----------------\n"
	for r := 7; r >= 0; r-- {
		for c := 0; c < 8; c++ {
			retval += fmt.Sprintf("|%c", pieces[r*8+c])
		}
		retval += "|\n-----------------\n"
	}
	return retval
}

func addPiece(bitmap uint64, charCode byte, pieces []byte) {
	for i := uint64(0); i < 64; i++ {
		if bitmap&(1<<i) != 0 {
			pieces[i] = charCode
		}
	}
}

type Player interface {
	GetNextMove(ctx context.Context, board *Board) (*Move, error)
	GetTimePerTurn() time.Duration
}

type ChessGame struct {
	whitePlayer Player
	blackPlayer Player
	board       *Board
	turnIndex   ColorIndex
}

func (b *Board) IsValidMove(m *Move, colorIndex ColorIndex) bool {
	return true
}

func (b *Board) applyMove(m *Move) error {
	return nil
}

func (g *ChessGame) IsFinished() bool {
	return false
}

func (g *ChessGame) GetWinner() Player {
	return nil
}

func NewStandardChessGame() *ChessGame {
	g := &ChessGame{}
	g.board = NewStandardBoard()
	g.whitePlayer = &HumanPlayer{}
	g.blackPlayer = &HumanPlayer{}
	g.turnIndex = WhiteIndex

	return g
}

func (g *ChessGame) PlayGame(ctx context.Context) {

	var err error
	for !g.IsFinished() {
		if g.turnIndex == BlackIndex {
			err = g.TakeTurn(ctx, g.blackPlayer)
		} else {
			err = g.TakeTurn(ctx, g.whitePlayer)
		}

		if err != nil {
			fmt.Printf("error taking turn, trying again. err: %s\n", err)
		}
		fmt.Printf("board:\n%s\n", g.board)
	}

}

func (g *ChessGame) updateTurnIndex() {
	if g.turnIndex == WhiteIndex {
		g.turnIndex = BlackIndex
	} else {
		g.turnIndex = WhiteIndex
	}
}

func (g *ChessGame) TakeTurn(ctx context.Context, p Player) error {
	pCtx, cancel := context.WithTimeout(ctx, p.GetTimePerTurn())
	defer cancel()
	m, err := p.GetNextMove(pCtx, g.board)
	if err != nil {
		return err
	}
	if !g.board.IsValidMove(m, g.turnIndex) {
		return fmt.Errorf("move not valid. Move %s", m)
	}
	err = g.board.applyMove(m)
	if err != nil {
		return err
	}
	g.updateTurnIndex()
	return nil
}
