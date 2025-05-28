package main

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
func Copy(g *GameState) *GameState {
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

	// Generateepassant moves page 31

}
