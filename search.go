package main

import (
	"fmt"
	"math"
	"time"
)

// result of a search
type SearchResult struct {
	BestMove     Move
	Score        int
	Depth        int
	NodesVisited int
	Duration     time.Duration
}

// Engine represents the chess engine
type Engine struct {
	MaxDepth     int
	TimeLimit    time.Duration
	NodesVisited int
	StartTime    time.Time
}

// creates new chess engine
func NewEngine() *Engine {
	return &Engine{
		MaxDepth:  5,
		TimeLimit: 5 * time.Second,
	}
}

// Minimax implements the minimax algorithm
func (e *Engine) Minimax(game *GameState, depth int, maximizingPlayer bool) int {
	e.NodesVisited++

	// Check time limit
	if time.Since(e.StartTime) > e.TimeLimit {
		return game.EvaluatePosition()
	}

	// Base case: maximum depth reached or game over
	if depth == 0 {
		return game.EvaluatePosition()
	}

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		// Game over
		if game.Board.IsInCheck(game.CurrentPlayer) {
			// Checkmate - return a score based on depth to prefer quicker mates
			if maximizingPlayer {
				return -30000 + depth
			}
			return 30000 - depth
		}
		// Stalemate
		return 0
	}

	if maximizingPlayer {
		maxEval := math.MinInt32
		for _, move := range moves {
			// Make the move
			newGame := game.Copy()
			newGame.MakeMove(move)

			// Recursive call
			eval := e.Minimax(newGame, depth-1, false)
			maxEval = max(maxEval, eval)
		}
		return maxEval
	} else {
		minEval := math.MaxInt32
		for _, move := range moves {
			// Make the move
			newGame := game.Copy()
			newGame.MakeMove(move)

			// Recursive call
			eval := e.Minimax(newGame, depth-1, true)
			minEval = min(minEval, eval)
		}
		return minEval
	}
}

// SearchBestMove finds the best move using minimax
func (e *Engine) SearchBestMove(game *GameState) SearchResult {
	e.StartTime = time.Now()
	e.NodesVisited = 0

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		return SearchResult{}
	}

	bestMove := moves[0]
	bestScore := math.MinInt32

	// Determine if we're maximizing (White) or minimizing (Black)
	maximizing := game.CurrentPlayer == White
	if !maximizing {
		bestScore = math.MaxInt32
	}

	fmt.Printf("Searching at depth %d...\n", e.MaxDepth)

	for i, move := range moves {
		// Make the move
		newGame := game.Copy()
		newGame.MakeMove(move)

		// Search
		score := e.Minimax(newGame, e.MaxDepth-1, !maximizing)

		fmt.Printf("Move %d/%d: %s -> %+d\n", i+1, len(moves), move.String(), score)

		// Check if this is the best move
		if maximizing {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		} else {
			if score < bestScore {
				bestScore = score
				bestMove = move
			}
		}

		// Check time limit
		if time.Since(e.StartTime) > e.TimeLimit {
			fmt.Println("Time limit reached!")
			break
		}
	}

	duration := time.Since(e.StartTime)

	return SearchResult{
		BestMove:     bestMove,
		Score:        bestScore,
		Depth:        e.MaxDepth,
		NodesVisited: e.NodesVisited,
		Duration:     duration,
	}
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Search with alpha beta pruning
func (e *Engine) AlphaBeta(game *GameState, depth int, alpha, beta int, maximizingPlayer bool) int {
	e.NodesVisited++

	// Check time limit
	if time.Since(e.StartTime) > e.TimeLimit {
		return game.EvaluatePosition()
	}

	// Base case
	if depth == 0 {
		return game.EvaluatePosition()
	}

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		// Game over
		if game.Board.IsInCheck(game.CurrentPlayer) {
			if maximizingPlayer {
				return -30000 + depth
			}
			return 30000 - depth
		}
		return 0
	}

	if maximizingPlayer {
		maxEval := math.MinInt32
		for _, move := range moves {
			newGame := game.Copy()
			newGame.MakeMove(move)

			eval := e.AlphaBeta(newGame, depth-1, alpha, beta, false)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)

			// Beta cutoff
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := math.MaxInt32
		for _, move := range moves {
			newGame := game.Copy()
			newGame.MakeMove(move)

			eval := e.AlphaBeta(newGame, depth-1, alpha, beta, true)
			minEval = min(minEval, eval)
			beta = min(beta, eval)

			// Alpha cutoff
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

// SearchBestMoveAB finds the best move using alpha-beta pruning
func (e *Engine) SearchBestMoveAB(game *GameState) SearchResult {
	e.StartTime = time.Now()
	e.NodesVisited = 0

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		return SearchResult{}
	}

	bestMove := moves[0]
	bestScore := math.MinInt32

	maximizing := game.CurrentPlayer == White
	if !maximizing {
		bestScore = math.MaxInt32
	}

	fmt.Printf("Searching with alpha-beta at depth %d...\n", e.MaxDepth)

	alpha := math.MinInt32
	beta := math.MaxInt32

	for i, move := range moves {
		newGame := game.Copy()
		newGame.MakeMove(move)

		score := e.AlphaBeta(newGame, e.MaxDepth-1, alpha, beta, !maximizing)

		fmt.Printf("Move %d/%d: %s -> %+d\n", i+1, len(moves), move.String(), score)

		if maximizing {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
			alpha = max(alpha, score)
		} else {
			if score < bestScore {
				bestScore = score
				bestMove = move
			}
			beta = min(beta, score)
		}

		if time.Since(e.StartTime) > e.TimeLimit {
			fmt.Println("Time limit reached!")
			break
		}
	}

	duration := time.Since(e.StartTime)

	return SearchResult{
		BestMove:     bestMove,
		Score:        bestScore,
		Depth:        e.MaxDepth,
		NodesVisited: e.NodesVisited,
		Duration:     duration,
	}
}

// ScoreMove assigns a score to a move for ordering purposes
func (e *Engine) ScoreMove(game *GameState, move Move) int {
	score := 0

	// Prioritize captures
	if move.IsCapture {
		// MVV-LVA (Most Valuable Victim - Least Valuable Attacker)
		victimValue := PieceValues[move.CapturedPiece.Type]
		attackerValue := PieceValues[move.PieceType]
		score += victimValue - attackerValue/10
	}

	// Prioritize promotions
	if move.PromotionPiece != Empty {
		score += PieceValues[move.PromotionPiece]
	}

	// Prioritize checks
	newGame := game.Copy()
	newGame.makeMove(move)
	enemyColor := 1 - game.CurrentPlayer
	if newGame.Board.IsInCheck(enemyColor) {
		score += 50
	}

	// Prioritize castling
	if move.IsCastle {
		score += 40
	}

	// Penalize moving to attacked squares
	enemyMoves := game.Board.GenerateAllMoves(1 - game.CurrentPlayer)
	for _, enemyMove := range enemyMoves {
		if enemyMove.ToRow == move.ToRow && enemyMove.ToCol == move.ToCol {
			score -= 10
			break
		}
	}

	return score
}

// OrderMoves sorts moves by their estimated value
func (e *Engine) OrderMoves(game *GameState, moves []Move) []Move {
	// Create a slice of move-score pairs
	type MoveScore struct {
		Move  Move
		Score int
	}

	moveScores := make([]MoveScore, len(moves))
	for i, move := range moves {
		moveScores[i] = MoveScore{
			Move:  move,
			Score: e.ScoreMove(game, move),
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(moveScores)-1; i++ {
		for j := i + 1; j < len(moveScores); j++ {
			if moveScores[j].Score > moveScores[i].Score {
				moveScores[i], moveScores[j] = moveScores[j], moveScores[i]
			}
		}
	}

	// Extract sorted moves
	orderedMoves := make([]Move, len(moves))
	for i, ms := range moveScores {
		orderedMoves[i] = ms.Move
	}

	return orderedMoves
}

// AlphaBetaOrdered implements alpha-beta with move ordering
func (e *Engine) AlphaBetaOrdered(game *GameState, depth int, alpha, beta int, maximizingPlayer bool) int {
	e.NodesVisited++

	if time.Since(e.StartTime) > e.TimeLimit {
		return game.EvaluatePosition()
	}

	if depth == 0 {
		return game.EvaluatePosition()
	}

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		if game.Board.IsInCheck(game.CurrentPlayer) {
			if maximizingPlayer {
				return -30000 + depth
			}
			return 30000 - depth
		}
		return 0
	}

	// Order moves for better pruning
	moves = e.OrderMoves(game, moves)

	if maximizingPlayer {
		maxEval := math.MinInt32
		for _, move := range moves {
			newGame := game.Copy()
			newGame.MakeMove(move)

			eval := e.AlphaBetaOrdered(newGame, depth-1, alpha, beta, false)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)

			if beta <= alpha {
				break // Beta cutoff
			}
		}
		return maxEval
	} else {
		minEval := math.MaxInt32
		for _, move := range moves {
			newGame := game.Copy()
			newGame.MakeMove(move)

			eval := e.AlphaBetaOrdered(newGame, depth-1, alpha, beta, true)
			minEval = min(minEval, eval)
			beta = min(beta, eval)

			if beta <= alpha {
				break // Alpha cutoff
			}
		}
		return minEval
	}
}

// SearchBestMoveOrdered finds the best move using ordered alpha-beta
func (e *Engine) SearchBestMoveOrdered(game *GameState) SearchResult {
	e.StartTime = time.Now()
	e.NodesVisited = 0

	moves := game.GenerateAllLegalMoves()
	if len(moves) == 0 {
		return SearchResult{}
	}

	// Order moves at root level
	moves = e.OrderMoves(game, moves)

	bestMove := moves[0]
	bestScore := math.MinInt32

	maximizing := game.CurrentPlayer == White
	if !maximizing {
		bestScore = math.MaxInt32
	}

	fmt.Printf("Searching with ordered alpha-beta at depth %d...\n", e.MaxDepth)

	alpha := math.MinInt32
	beta := math.MaxInt32

	for i, move := range moves {
		newGame := game.Copy()
		newGame.MakeMove(move)

		score := e.AlphaBetaOrdered(newGame, e.MaxDepth-1, alpha, beta, !maximizing)

		fmt.Printf("Move %d/%d: %s -> %+d\n", i+1, len(moves), move.String(), score)

		if maximizing {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
			alpha = max(alpha, score)
		} else {
			if score < bestScore {
				bestScore = score
				bestMove = move
			}
			beta = min(beta, score)
		}

		if time.Since(e.StartTime) > e.TimeLimit {
			fmt.Println("Time limit reached!")
			break
		}
	}

	duration := time.Since(e.StartTime)

	return SearchResult{
		BestMove:     bestMove,
		Score:        bestScore,
		Depth:        e.MaxDepth,
		NodesVisited: e.NodesVisited,
		Duration:     duration,
	}
}
