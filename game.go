package main

import "strings"

// GameState holds the state of the chess match
type GameState struct {
	Board           *Board
	CurrentPlayer   int
	MoveHistory     []Move
	WhiteCanCastleK bool // Kingside castling
	WhiteCanCastleQ bool // Queenside castling
	BlackCanCastleK bool
	BlackCanCastleQ bool
	EnPassantSquare [2]int // Row , Col of EnPassant target square (-1, -1) if none
	HalfMoveClock   int    // For 50 move rule
	FullMoveNumber  int
}

// Creates new chess game
func NewGame() *GameState {
	return &GameState{
		Board:           NewBoard(),
		CurrentPlayer:   White,
		MoveHistory:     make([]Move, 0),
		WhiteCanCastleK: true,
		WhiteCanCastleQ: true,
		BlackCanCastleK: true,
		BlackCanCastleQ: true,
		EnPassantSquare: [2]int{-1, -1},
		HalfMoveClock:   0,
		FullMoveNumber:  1,
	}
}

// Create copy of game state (deep copy)
func (g *GameState) Copy() *GameState {
	newGame := &GameState{
		Board:           g.Board.Copy(),
		CurrentPlayer:   g.CurrentPlayer,
		MoveHistory:     make([]Move, len(g.MoveHistory)),
		WhiteCanCastleK: g.WhiteCanCastleK,
		WhiteCanCastleQ: g.WhiteCanCastleQ,
		BlackCanCastleK: g.BlackCanCastleK,
		BlackCanCastleQ: g.BlackCanCastleQ,
		EnPassantSquare: g.EnPassantSquare,
		HalfMoveClock:   g.HalfMoveClock,
		FullMoveNumber:  g.FullMoveNumber,
	}
	copy(newGame.MoveHistory, g.MoveHistory) // What does this do, is it an inbuilt array function?
	return newGame
}

// Check to see whether castling is possible?
func (g *GameState) CanCastle(color, side int) bool {
	// side: 0 = kingside, 1 = queenside

	// check castling rights
	var canCastle bool
	if color == White {
		canCastle = (side == 0 && g.WhiteCanCastleK) || (side == 1 && g.WhiteCanCastleQ)
	} else {
		canCastle = (side == 0 && g.BlackCanCastleK) || (side == 1 && g.BlackCanCastleQ)
	}

	if !canCastle {
		return false
	}

	// King must not be in check
	if g.Board.IsInCheck(color) {
		return false
	}

	// White black rank
	row := 7
	if color == Black {
		row = 0
	}

	var squares []int
	if side == 0 { // Kingside
		squares = []int{5, 6} //f1 f8, g1 g8
	} else { // Queenside
		squares = []int{1, 2, 3} // b1 b8, c1 c8, d1 d8
	}

	// check if squares are empty
	for _, col := range squares {
		if g.Board.GetPiece(row, col).Type != Empty {
			return false
		}
	}

	// King cannot pass through or end up in check
	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	checkSquares := []int{4, 5, 6} // Kings path for kingside
	if side == 1 {                 // Queenside
		checkSquares = []int{4, 3, 2}
	}

	for _, col := range checkSquares {
		if g.Board.IsSquareAttacked(row, col, enemyColor) {
			return false
		}
	}

	return true

}

// Adding castling moves to move list

func (g *GameState) GenerateCastlingMoves(moves *[]Move) {
	color := g.CurrentPlayer

	// Kingside castling
	if g.CanCastle(color, 0) {
		row := 7
		if color == Black {
			row = 0
		}

		move := Move{
			FromRow: row, FromCol: 4,
			ToRow: row, ToCol: 6,
			PieceType: King,
			IsCastle:  true,
		}

		*moves = append(*moves, move)

	}

	// Queenside castling
	if g.CanCastle(color, 1) {
		row := 7
		if color == Black {
			row = 0
		}

		move := Move{
			FromRow: row, FromCol: 4,
			ToRow: row, ToCol: 2,
			PieceType: King,
			IsCastle:  true,
		}

		*moves = append(*moves, move)
	}

}

func (g *GameState) GenerateEnPassantMoves(moves *[]Move) {
	if g.EnPassantSquare[0] == -1 {
		return // No EnPassant possible
	}

	targetRow := g.EnPassantSquare[0]
	targetCol := g.EnPassantSquare[1]
	color := g.CurrentPlayer

	// Find pawns that can capture enpassant
	pawnRow := targetRow + 1
	if color == Black {
		pawnRow = targetRow - 1
	}

	// Check left and right adjacent squares
	for _, deltaCol := range []int{-1, 1} {
		pawnCol := targetCol + deltaCol
		if IsValidSquare(pawnRow, pawnCol) {
			piece := g.Board.GetPiece(pawnRow, pawnCol)
			if piece.Type == Pawn && piece.Color == color {
				move := Move{FromRow: pawnRow, FromCol: pawnCol, ToRow: targetRow, ToCol: targetCol,
					PieceType: Pawn, IsEnPassant: true, IsCapture: true,
				}

				*moves = append(*moves, move)
			}
		}
	}

}

// Generate all legal moves including special moves

func (g *GameState) GenerateAllLegalMoves() []Move {
	moves := g.Board.GenerateAllMoves(g.CurrentPlayer)

	// Adding castling moves
	g.GenerateCastlingMoves(&moves) // can you explain how this works. It just appends the moves

	// Adding en passant moves
	g.GenerateEnPassantMoves(&moves)

	// filter legal moves
	var legalMoves []Move
	for _, move := range moves {
		if g.IsLegalMove(move) { // check this function
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

// Legal move from gamestate

func (g *GameState) IsLegalMove(move Move) bool {
	// Make a copy of game and try the move

	testGame := g.Copy()
	testGame.makeMove(move)

	return !testGame.Board.IsInCheck(g.CurrentPlayer)
}

// Make move executes move and updates game state - helper function
func (g *GameState) makeMove(move Move) {
	// Handle castling
	if move.IsCastle {
		// Move king
		g.Board.MakeMove(move)

		// Move rook
		row := move.FromRow
		if move.ToCol == 6 { // Kingside
			rookMove := Move{
				FromRow: row, FromCol: 7,
				ToRow: row, ToCol: 5,
				PieceType: Rook,
			}
			g.Board.MakeMove(rookMove)
		} else {
			// Queenside castling
			rookMove := Move{
				FromRow: row, FromCol: 0,
				ToRow: row, ToCol: 3,
				PieceType: Rook,
			}
			g.Board.MakeMove(rookMove)
		}

		return

	}
	if move.IsEnPassant {
		g.Board.MakeMove(move)

		// Remove Captured Pawn
		capturedRow := move.ToRow + 1
		if g.CurrentPlayer == Black {
			capturedRow = move.ToRow - 1
		}
		g.Board.SetPiece(capturedRow, capturedRow, Piece{Empty, White})

		return
	}

	// Regular move

	g.Board.MakeMove(move)
}

// Make move exceutes move and updates whole game state
func (g *GameState) MakeMove(move Move) bool {
	// Verify move is legal
	legalMoves := g.GenerateAllLegalMoves()
	isLegal := false

	for _, legalMove := range legalMoves {
		if move.FromRow == legalMove.FromRow && move.FromCol == legalMove.FromCol &&
			move.ToRow == legalMove.ToRow && move.ToCol == legalMove.ToCol {
			move = legalMove
			isLegal = true
			break
		}
	}

	if !isLegal {
		return false
	}

	// Update castling rights
	if move.PieceType == King {
		if g.CurrentPlayer == White {
			g.WhiteCanCastleK = false
			g.WhiteCanCastleQ = false
		} else {
			g.BlackCanCastleK = false
			g.BlackCanCastleQ = false
		}

	}

	if move.PieceType == Rook {
		if g.CurrentPlayer == White {
			if move.FromRow == 7 && move.FromCol == 0 {
				g.WhiteCanCastleQ = false
			} else if move.FromRow == 7 && move.FromCol == 7 {
				g.WhiteCanCastleK = false
			}
		} else {
			if move.FromRow == 0 && move.FromCol == 0 {
				g.BlackCanCastleQ = false
			} else if move.FromRow == 0 && move.FromCol == 7 {
				g.BlackCanCastleK = false
			}
		}
	}

	// Update Enpassant square
	g.EnPassantSquare = [2]int{-1, -1}

	if move.PieceType == Pawn && abs(move.ToRow-move.ToCol) == 2 {
		// double pawn move -> set enpassant square
		g.EnPassantSquare[0] = (move.FromRow - move.ToRow) / 2
		g.EnPassantSquare[1] = move.FromCol

	}

	// Update half move clock
	if move.PieceType == Pawn || move.IsCapture {
		g.HalfMoveClock = 0
	} else {
		g.HalfMoveClock++
	}

	// Execute move
	g.makeMove(move)

	// Update move history
	g.MoveHistory = append(g.MoveHistory, move)

	// Update current player
	g.CurrentPlayer = 1 - g.CurrentPlayer

	// update full move number
	if g.CurrentPlayer == White {
		g.FullMoveNumber++
	} // white starts so increment move each time it's white

	return true

}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// Is Over checks if game has ended

func (g *GameState) IsGameOver() (bool, string) {
	legalMoves := g.GenerateAllLegalMoves()

	if len(legalMoves) == 0 {
		if g.Board.IsInCheck(g.CurrentPlayer) {
			// Checkmate
			winner := "White"
			if g.CurrentPlayer == White {
				winner = "Black"
			}
			return true, winner + "Wins by checkmate!"
		} else {
			// Stalemate
			return true, "Draw by stalemate."
		}
	}

	// 50 mve rule
	if g.HalfMoveClock >= 100 { // 50 moves each
		return true, "Draw by 50 move rule"
	}

	return false, ""

}

// Return current player as string

func (g *GameState) GetCurrentPlayerString() string {
	if g.CurrentPlayer == White {
		return "White"
	}

	return "Black"
}

// Parse move converts algebraic notation to a Move struct
func (g *GameState) ParseMove(notation string) (Move, bool) {
	notation = strings.TrimSpace(notation)

	if notation == "O-O" || notation == "0-0" {
		// kingside castling
		row := 7
		if g.CurrentPlayer == Black {
			row = 0
		}

		return Move{
			FromRow: row, FromCol: 4,
			ToRow: row, ToCol: 6,
			PieceType: King,
			IsCastle:  true,
		}, true
	}
	if notation == "O-O-O" || notation == "0-0-0" {
		// Queenside castling
		row := 7
		if g.CurrentPlayer == Black {
			row = 0
		}

		return Move{
			FromRow: row, FromCol: 4,
			ToRow: row, ToCol: 2,
			PieceType: King,
			IsCastle:  true,
		}, true
	}

	// Simple format e2e4

	if len(notation) >= 4 {
		files := "abcdefgh"

		fromCol := strings.IndexByte(files, notation[0])
		fromRow := 8 - int(notation[1]-'O') // What does this mean?
		toCol := strings.IndexByte(files, notation[2])
		toRow := 8 - int(notation[3]-'O')

		if fromCol >= 0 && fromRow >= 0 && fromRow < 8 && toCol >= 0 &&
			toRow >= 0 && toRow < 8 {
			piece := g.Board.GetPiece(fromRow, fromCol)
			move := Move{
				FromRow: fromRow, FromCol: fromCol,
				ToRow: toRow, ToCol: toCol,
				PieceType: piece.Type,
			}

			// Check promotion

			if len(notation) >= 5 && piece.Type == Pawn {
				switch notation[4] {
				case 'q', 'Q':
					move.PromotionPiece = Queen
				case 'r', 'R':
					move.PromotionPiece = Rook
				case 'n', 'N':
					move.PromotionPiece = Knight
				case 'b', 'B':
					move.PromotionPiece = Bishop
				}
			}

			return move, true
		}

	}

	return Move{}, false

}
