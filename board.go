package main

import (
	"fmt"
)

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
			Empty: ".", Pawn: "P", Rook: "R", Knight: "N", Bishop: "B", Queen: "Q", King: "K",
		},
		Black: {
			Empty: ".", Pawn: "p", Rook: "r", Knight: "n", Bishop: "b", Queen: "q", King: "k",
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

func (b *Board) GenerateRookMoves(row, col int, moves *[]Move) {
	piece := b.GetPiece(row, col)
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // Up down left right

	for _, dir := range directions {
		for distance := 1; distance < 8; distance++ {
			newRow := row + dir[0]*distance
			newCol := col + dir[1]*distance

			if !IsValidSquare(newRow, newCol) {
				break
			}

			target := b.GetPiece(newRow, newCol)
			if target.Type == Empty {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType: Rook,
				}
				*moves = append(*moves, move)
			} else if target.Color != piece.Color {
				// Oppositon piece, can capture
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType:     Rook,
					CapturedPiece: target,
					IsCapture:     true,
				}
				*moves = append(*moves, move)
				break // cant move any further in this direction
			} else {
				// Own piece can't move
				break
			}
		}
	}
}

// Bishop moves

func (b *Board) GenerateBishopMoves(row, col int, moves *[]Move) {
	piece := b.GetPiece(row, col)
	directions := [][]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}} // Diagonals

	for _, dir := range directions {
		for distance := 1; distance < 8; distance++ {
			newRow := row + dir[0]*distance
			newCol := col + dir[1]*distance

			if !IsValidSquare(newRow, newCol) {
				break
			}

			target := b.GetPiece(newRow, newCol)
			if target.Type == Empty {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType: Bishop,
				}
				*moves = append(*moves, move)
			} else if target.Color != piece.Color {
				// capture
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType:     Bishop,
					CapturedPiece: target,
					IsCapture:     true,
				}
				*moves = append(*moves, move)
				break
			} else {
				break
			}

		}

	}
}

// Queen moves

func (b *Board) GenerateQueenMoves(row, col int, moves *[]Move) {
	b.GenerateBishopMoves(row, col, moves) // Queen moves like both a rook and bishop
	b.GenerateRookMoves(row, col, moves)
}

// Knight moves
func (b *Board) GenerateKnightMoves(row, col int, moves *[]Move) {
	piece := b.GetPiece(row, col)
	knightMoves := [][]int{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1},
	}

	for _, delta := range knightMoves {
		newRow := row + delta[0]
		newCol := col + delta[1]

		if IsValidSquare(newRow, newCol) {
			target := b.GetPiece(newRow, newCol)
			if target.Type == Empty {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType: Knight,
				}
				*moves = append(*moves, move)
			} else if target.Color != piece.Color {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType:     Knight,
					CapturedPiece: target,
					IsCapture:     true,
				}
				*moves = append(*moves, move)

			}
		}

	}
}

// King moves

func (b *Board) GenerateKingMoves(row, col int, moves *[]Move) {
	piece := b.GetPiece(row, col)
	kingMoves := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1}, {1, -1}, {1, 0},
		{1, 1},
	}
	for _, delta := range kingMoves {
		newRow := row + delta[0]
		newCol := col + delta[1]

		if IsValidSquare(row, col) {
			target := b.GetPiece(newRow, newCol)
			if target.Type == Empty {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType: King,
				}
				*moves = append(*moves, move)
			} else if target.Color != piece.Color {
				move := Move{
					FromRow: row, FromCol: col,
					ToRow: newRow, ToCol: newCol,
					PieceType:     King,
					CapturedPiece: target,
					IsCapture:     true,
				}
				*moves = append(*moves, move)
			}
		}
	}
}

// Generate all moves

func (b *Board) GenerateAllMoves(color int) []Move {
	var moves []Move

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := b.GetPiece(row, col)

			if piece.Type != Empty && piece.Color == color {
				switch piece.Type {
				case Pawn:
					b.GeneratePawnMoves(row, col, &moves)
				case Rook:
					b.GenerateRookMoves(row, col, &moves)
				case Bishop:
					b.GenerateBishopMoves(row, col, &moves)
				case Knight:
					b.GenerateKnightMoves(row, col, &moves)
				case Queen:
					b.GenerateQueenMoves(row, col, &moves)
				case King:
					b.GenerateKingMoves(row, col, &moves)
				}
			}
		}

	}
	return moves
}

// Find King
func (b *Board) FindKing(color int) (int, int) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := b.GetPiece(row, col)
			if piece.Type == King && piece.Color == color {
				return row, col
			}
		}
	}
	return -1, -1 // Should never happen
}

// Is Square attacked -0 checks if square is under attack by given color

func (b *Board) IsSquareAttacked(row, col, attackingColor int) bool {
	// Check for pawn attacks
	pawnDirection := 1

	if attackingColor == Black {
		pawnDirection = -1
	}
	// Check diagonal pawn attack
	for _, deltaCol := range []int{-1, 1} {
		pawnRow := row + pawnDirection
		pawnCol := col + deltaCol
		if IsValidSquare(pawnRow, pawnCol) {
			piece := b.GetPiece(pawnRow, pawnCol)
			if piece.Type == Pawn && piece.Color == attackingColor {
				return true
			}
		}
	}
	// Check Knight attacks
	knightMoves := [][]int{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1},
	}

	for _, delta := range knightMoves {
		knightRow := row + delta[0]
		knightCol := col + delta[1]
		if IsValidSquare(knightRow, knightCol) {
			piece := b.GetPiece(knightRow, knightCol)
			if piece.Type == Knight && piece.Color == attackingColor {
				return true
			}
		}
	}

	// Check sliding attacks rook bishop queen
	directions := [][]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}

	for i, dir := range directions {
		for distance := 1; distance < 8; distance++ {
			checkRow := row + dir[0]
			checkCol := col + dir[1]

			if IsValidSquare(checkRow, checkCol) {
				break
			}

			piece := b.GetPiece(checkRow, checkCol)
			if piece.Type != Empty {
				if piece.Color == attackingColor {
					if i < 4 { // rook directions
						if piece.Type == Rook || piece.Type == Queen {
							return true
						} else { // bishop directions
							if piece.Type == Bishop || piece.Type == Queen {
								return true
							}
						}
					}
					break // piece blocks further movement in direction
				}

			}
		}
	}
	kingMoves := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1}, {1, -1},
		{1, 0}, {1, 1},
	}

	for _, delta := range kingMoves {
		kingRow := row + delta[0]
		kingCol := col + delta[1]

		if IsValidSquare(kingRow, kingCol) {
			piece := b.GetPiece(kingRow, kingCol)
			if piece.Type == King && piece.Color == attackingColor {
				return true
			}
		}
	}
	return false
}

// Is in check cheks whether that colors king is in check

func (b *Board) IsInCheck(color int) bool {
	kingRow, kingCol := b.FindKing(color)
	if kingRow == -1 {
		return false // No king found - should never happen
	}

	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	return b.IsSquareAttacked(kingRow, kingCol, enemyColor)
}

// Copy board - deep copy of game

func (b *Board) Copy() *Board {
	newBoard := &Board{}
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			newBoard.squares[row][col] = b.squares[row][col]
		}
	}
	return newBoard
}

func (b *Board) IsLegalMove(move Move, color int) bool {
	testBoard := b.Copy()
	testBoard.MakeMove(move)

	return !testBoard.IsInCheck(color)
}

// Generate all legal moves for color

func (b *Board) GenerateLegalMoves(color int) []Move {
	pseudoLegalMoves := b.GenerateAllMoves(color)
	var legalMoves []Move
	for _, move := range pseudoLegalMoves {
		if b.IsLegalMove(move, color) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

// Game state detection
// Is checkmate
func (b *Board) IsCheckmate(color int) bool {
	if !b.IsInCheck(color) {
		return false
	}

	legalMoves := b.GenerateLegalMoves(color)
	return len(legalMoves) == 0 // If true, then no available moves to player
}

func (b *Board) IsStaleMate(color int) bool {
	if b.IsInCheck(color) {
		return false
	}
	legalMoves := b.GenerateLegalMoves(color)
	return len(legalMoves) == 0

}

// Get game result

func (b *Board) GetGameResult(currentPlayer int) string {
	if b.IsCheckmate(currentPlayer) {
		if currentPlayer == White {
			return "Black wins by checkmate"
		}
		return "White wins by checkmate"
	}

	if b.IsStaleMate(currentPlayer) {
		return "Draw by stalemate"
	}

	if b.IsInCheck(currentPlayer) {
		return "Check"
	}

	return "Game continues!"
}
