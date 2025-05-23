package main

import "fmt"

// Pieces
const (
	Empty = iota
	Pawn
	Rook
	Bishop
	Knight
	Queen
	King
)

// Colours - whats iota?
const (
	White = iota
	Black
)

type Piece struct {
	Type  int
	Color int
}

type Board struct {
	squares [8][8]Piece
}

// New board function
func NewBoard() *Board {
	b := &Board{}
	b.setupStartingPosition()
	return b
}

func (b *Board) setupStartingPosition() {
	// Clear board
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			b.squares[row][col] = Piece{Empty, White}
		}
	}

	// Set up white pieces
	backRank := []int{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for col := 0; col < 8; col++ {
		b.squares[7][col] = Piece{backRank[col], White}
		b.squares[6][col] = Piece{Pawn, White}
	}

	// Setting up black pieces
	for col := 0; col < 8; col++ {
		b.squares[0][col] = Piece{backRank[col], Black}
		b.squares[1][col] = Piece{Pawn, Black}

	}
}

// Get piece to return piece from that position
func (b *Board) GetPiece(row, col int) Piece {
	if row < 0 || row >= 8 || col < 0 || col >= 8 {
		return Piece{Empty, White}
	} else {
		return b.squares[row][col]
	}
}

// Set Piece sets the piece at the given position
func (b *Board) SetPiece(row, col int, piece Piece) {
	if row >= 0 && row < 8 && col >= 0 && col < 8 {
		b.squares[row][col] = piece
	}
}

// Display board to print board to terminal
func (b *Board) Display() {
	pieceSymbols := map[int]map[int]string{
		White: {
			Empty: ".", Pawn: "P", Rook: "R", Knight: "K", Bishop: "B", Queen: "Q", King: "K",
		},
		Black: {
			Empty: ".", Pawn: "p", Rook: "r", Knight: "k", Bishop: "b", Queen: "q", King: "k",
		},
	}
	fmt.Println(" a b c d e f g h")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d", 8-row)
		for col := 0; col < 8; col++ {
			piece := b.squares[row][col]
			symbol := pieceSymbols[piece.Color][piece.Type]
			fmt.Printf("%s ", symbol)
		}
		fmt.Printf("%d\n", 8-row)
	}

	fmt.Println(" a b c d e f g h")

}
