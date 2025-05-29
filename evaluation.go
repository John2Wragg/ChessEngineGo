package main

import "math"

var PieceValues = map[int]int{
	Empty:  0,
	Pawn:   100,
	Knight: 320,
	Bishop: 330,
	Rook:   500,
	Queen:  900,
	King:   20000, // Invaluable
}

// piece square tables for evaluation
var PawnTable = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{50, 50, 50, 50, 50, 50, 50, 50},
	{10, 10, 20, 30, 30, 20, 10, 10},
	{5, 5, 10, 25, 25, 10, 5, 5},
	{0, 0, 0, 20, 20, 0, 0, 0},
	{5, -5, -10, 0, 0, -10, -5, 5},
	{5, 10, 10, -20, -20, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var KnightTable = [8][8]int{
	{-50, -40, -30, -30, -30, -30, -40, -50},
	{-40, -20, 0, 0, 0, 0, -20, -40},
	{-30, 0, 10, 15, 15, 10, 0, -30},
	{-30, 5, 15, 20, 20, 15, 5, -30},
	{-30, 0, 15, 20, 20, 15, 0, -30},
	{-30, 5, 10, 15, 15, 10, 5, -30},
	{-40, -20, 0, 5, 5, 0, -20, -40},
	{-50, -40, -30, -30, -30, -30, -40, -50},
}

var BishopTable = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 5, 5, 10, 10, 5, 5, -10},
	{-10, 0, 10, 10, 10, 10, 0, -10},
	{-10, 10, 10, 10, 10, 10, 10, -10},
	{-10, 5, 0, 0, 0, 0, 5, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

var RookTable = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{0, 0, 0, 5, 5, 0, 0, 0},
}

var QueenTable = [8][8]int{
	{-20, -10, -10, -5, -5, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 5, 5, 5, 0, -10},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-10, 5, 5, 5, 5, 5, 0, -10},
	{-10, 0, 5, 0, 0, 0, 0, -10},
	{-20, -10, -10, -5, -5, -10, -10, -20},
}

var KingMiddleGameTable = [8][8]int{
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-10, -20, -20, -20, -20, -20, -20, -10},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{20, 30, 10, 0, 0, 10, 30, 20},
}

var KingEndGameTable = [8][8]int{
	{-50, -40, -30, -20, -20, -30, -40, -50},
	{-30, -20, -10, 0, 0, -10, -20, -30},
	{-30, -10, 20, 30, 30, 20, -10, -30},
	{-30, -10, 30, 40, 40, 30, -10, -30},
	{-30, -10, 30, 40, 40, 30, -10, -30},
	{-30, -10, 20, 30, 30, 20, -10, -30},
	{-30, -30, 0, 0, 0, 0, -30, -30},
	{-50, -30, -30, -30, -30, -30, -30, -50},
}

// Get piece square table
func GetPieceSquareTable(piceType int, isEndgame bool) [8][8]int {
	switch piceType {
	case Pawn:
		return PawnTable
	case Knight:
		return KnightTable
	case Bishop:
		return BishopTable
	case Rook:
		return RookTable
	case Queen:
		return QueenTable
	case King:
		if isEndgame {
			return KingEndGameTable
		}
		return KingMiddleGameTable
	default:
		return [8][8]int{}
	}
}

// Evaluate position
// Positive values for white, negative for black
func (g *GameState) EvaluatePosition() int {
	score := 0

	// Check for checkmate and stalemate
	moves := g.GenerateAllLegalMoves()
	if len(moves) == 0 {
		if g.Board.IsInCheck(g.CurrentPlayer) {
			if g.CurrentPlayer == White {
				return -30000
			}
			return 30000
		}

		return 0 // Stalemate
	}

	whiteMaterial, blackMaterial := 0, 0

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := g.Board.GetPiece(row, col)
			if piece.Type == Empty {
				continue
			}

			materialValue := PieceValues[piece.Type]

			// Positional Value
			isEndgame := g.isEndgame() // Why is this erroring?
			table := GetPieceSquareTable(piece.Type, isEndgame)

			var positionalValue int
			if piece.Color == White {
				positionalValue = table[row][col]
				whiteMaterial += materialValue
			} else {
				positionalValue = table[7-row][col]
				blackMaterial += materialValue
			}

			totalValue := materialValue + positionalValue

			if piece.Color == White {
				score += totalValue
			} else {
				score -= totalValue
			}

		}

	}

	// mobility bonus (number of legal moves)

	whiteMoves := len(g.Board.GenerateAllMoves(White))
	blackMoves := len(g.Board.GenerateAllMoves(Black))
	mobilityScore := (whiteMoves - blackMoves) * 10
	score += mobilityScore

	// Castling bonus
	if g.WhiteCanCastleK || g.WhiteCanCastleQ {
		score += 20
	}

	if g.BlackCanCastleK || g.BlackCanCastleQ {
		score -= 20
	}

	// King safety in middle game
	if !g.isEndgame() {
		score += g.evaluateKingSafety(White) - g.evaluateKingSafety(Black)
	}

	return score
}

// isEndgame determines whether we're in an end game
func (g *GameState) isEndgame() bool {
	pieceCount := 0
	queens := 0

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := g.Board.GetPiece(row, col)
			if piece.Type != Empty && piece.Type != King {
				pieceCount++
				if piece.Type == Queen {
					queens++
				}
			}

		}
	}
	// End game if no queens or very few pieces in total
	return (queens == 0 && pieceCount < 12) || (pieceCount < 8)
}

// Evaluate kingSafety returns a score for the safety
func (g *GameState) evaluateKingSafety(color int) int {
	kingRow, kingCol := g.Board.FindKing(color)

	if kingRow == -1 {
		return -1000 // No king found (shouldn't ever happen)
	}

	safety := 0
	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	// Check squares around king
	for deltaRow := -1; deltaRow <= 1; deltaRow++ {
		for deltaCol := -1; deltaCol <= 1; deltaCol++ {
			if deltaRow == 0 && deltaCol == 0 {
				continue
			}
			checkRow := kingRow + deltaRow
			checkCol := kingCol + deltaCol

			if IsValidSquare(checkRow, checkCol) {
				if g.Board.IsSquareAttacked(checkRow, checkCol, enemyColor) {
					safety -= 10
				}
				// Bonus if pieces are defending king
				piece := g.Board.GetPiece(checkRow, checkCol)
				if piece.Type != Empty && piece.Color == color {
					safety += 5
				}

			}

		}
	}

	return safety

}

// Get Game Phase return value from 0 - endgame, to 1 - opening
func (g *GameState) GetGamePhase() float64 {
	totalMaterial := 0

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := g.Board.GetPiece(row, col)
			if piece.Type != Empty && piece.Type != King {
				totalMaterial += PieceValues[piece.Type]
			}
		}
	}

	phase := float64(totalMaterial) / 7800.0 // max without kings
	return math.Min(1.0, math.Max(0.0, phase))

}
