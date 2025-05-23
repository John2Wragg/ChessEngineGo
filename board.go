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

// Adding moves
type Move struct {
	FromRow, FromCol int
	ToRow, ToCol     int
	PieceType        int
	CapturedPiece    Piece
	IsCapture        bool
	IsEnPassant      bool
	IsCastle         bool
	PromotionPiece   int
}

// String to return readable version of move
func (m Move) String() string {
	files := "abcdefgh"
	fromSquare := fmt.Sprintf("%c%d", files[m.FromCol], 8-m.FromRow)
	toSquare := fmt.Sprintf("%c%d", files[m.ToCol], 8-m.ToRow)
	return fromSquare + toSquare
}

// Make move executes a move on the board
func (b *Board) MakeMove(move Move) {
	piece := b.GetPiece(move.FromRow, move.FromCol)

	b.SetPiece(move.FromRow, move.FromCol, Piece{Empty, White})

	// Place piece at new destination
	if move.PromotionPiece != Empty {
		// Handle pawn promotion
		b.SetPiece(move.ToRow, move.ToCol, Piece{move.PromotionPiece, piece.Color})
	} else {
		b.SetPiece(move.ToRow, move.ToCol, piece)
	}
}

// Checks if coordinates are within board
func IsValidSquare(row, col int) bool {
	return row >= 0 && row < 8 && col >= 0 && col < 8
}

// Adding moves
// Generate Pawn Moves
func (b *Board) GeneratePawnMoves(row, col int, moves *[]Move) {
	piece := b.GetPiece(row, col)
	direction := -1 // white moves up the board (decreasing row numbers)
	startRow := 6

	if piece.Color == Black {
		direction = 1
		startRow = 1
	}

	// Forward move
	newRow := row + direction
	if IsValidSquare(newRow, col) && b.GetPiece(newRow, col).Type == Empty {
		move := Move{
			FromRow: row, FromCol: col,
			ToRow: newRow, ToCol: col,
			PieceType: Pawn,
		}

		// Checking for promotion
		if (piece.Color == White && newRow == 0) || (piece.Color == Black && newRow == 7) {
			promotionPieces := []int{Queen, Rook, Bishop, Knight}
			for _, promoPiece := range promotionPieces {
				promoMove := move
				promoMove.PromotionPiece = promoPiece
				*moves = append(*moves, promoMove)
			}
		} else {
			*moves = append(*moves, move)

			// if starting can double move
			if row == startRow {
				doubleRow := 2 * direction
				if IsValidSquare(doubleRow, col) && b.GetPiece(doubleRow, col).Type == Empty {
					doubleMove := Move{
						FromRow: row, FromCol: col,
						ToRow: doubleRow, ToCol: col,
						PieceType: Pawn,
					}
					*moves = append(*moves, doubleMove)
				}
			}
		}

	}

	// Capture moves
	for _, deltaCol := range []int{-1, 1} {
		newRow := row + direction
		newCol := col + deltaCol

		if IsValidSquare(newRow, newCol) {
			target := b.GetPiece(newRow, newCol)
			if target.Type != Empty && target.Color != piece.Color {
				move := Move{FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType:     Pawn,
					CapturedPiece: target,
					IsCapture:     true,
				}

				// Check for promotion
				if (piece.Color == White && newRow == 0) || (piece.Color == Black && newRow == 7) {
					promotionPieces := []int{Queen, Rook, Bishop, Knight}
					for _, promoPiece := range promotionPieces {
						promoMove := move
						promoMove.PromotionPiece = promoPiece
						*moves = append(*moves, promoMove)

					}

				} else {
					*moves = append(*moves, move)

				}
			}

		}

	}
}
