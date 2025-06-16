package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Chess Engine v1.0")
	fmt.Println("=================")
	fmt.Println("Commands: move, eval, ai, depth <n>, quit, moves, help")
	fmt.Println()

	game := NewGame()
	engine := NewEngine()
	scanner := bufio.NewScanner(os.Stdin)

	for {

		game.Board.Display()

		eval := game.EvaluatePosition()
		fmt.Printf("\nPosition evaluation: %+d centipawns", eval)
		if eval > 0 {
			fmt.Print(" (White is better)")
		} else if eval < 0 {
			fmt.Print(" (Black is better)")
		} else {
			fmt.Print(" (Equal position)")
		}
		fmt.Println()

		fmt.Printf("\n%s to move", game.GetCurrentPlayerString())

		// Check game status
		if gameOver, result := game.IsGameOver(); gameOver {
			fmt.Printf("\nGame Over: %s\n", result)
			break
		}

		if game.Board.IsInCheck(game.CurrentPlayer) {
			fmt.Print(" (in check)")
		}
		fmt.Print(": ")

		// Get input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "quit", "q":
			fmt.Println("Thanks for playing!")
			return

		case "eval", "e":
			eval := game.EvaluatePosition()
			phase := game.GetGamePhase()
			isEndgame := game.isEndgame()
			fmt.Printf("Detailed evaluation:\n")
			fmt.Printf("  Total score: %+d centipawns\n", eval)
			fmt.Printf("  Game phase: %.2f (1.0=opening, 0.0=endgame)\n", phase)
			fmt.Printf("  Is endgame: %t\n", isEndgame)

			// Count material
			whiteMaterial, blackMaterial := 0, 0
			for row := 0; row < 8; row++ {
				for col := 0; col < 8; col++ {
					piece := game.Board.GetPiece(row, col)
					if piece.Type != Empty {
						if piece.Color == White {
							whiteMaterial += PieceValues[piece.Type]
						} else {
							blackMaterial += PieceValues[piece.Type]
						}
					}
				}
			}
			fmt.Printf("  White material: %d\n", whiteMaterial)
			fmt.Printf("  Black material: %d\n", blackMaterial)
			fmt.Printf("  Material difference: %+d\n", whiteMaterial-blackMaterial)

		case "ai":
			fmt.Println("AI is thinking...")
			result := engine.SearchBestMoveOrdered(game)

			if result.BestMove.PieceType != Empty {
				fmt.Printf("\nAI chooses: %s (score: %+d)\n", result.BestMove.String(), result.Score)
				fmt.Printf("Search info: %d nodes, depth %d, %.2fs\n",
					result.NodesVisited, result.Depth, result.Duration.Seconds())

				if game.MakeMove(result.BestMove) {
					fmt.Println("Move played successfully")
				} else {
					fmt.Println("Error: AI suggested illegal move!")
				}
			} else {
				fmt.Println("AI couldn't find a move")
			}

		case "depth":
			if len(parts) >= 2 {
				if depth, err := strconv.Atoi(parts[1]); err == nil && depth > 0 && depth <= 10 {
					engine.MaxDepth = depth
					fmt.Printf("Search depth set to %d\n", depth)
				} else {
					fmt.Println("Invalid depth. Use 1-10")
				}
			} else {
				fmt.Printf("Current search depth: %d\n", engine.MaxDepth)
			}

		case "moves", "m":
			moves := game.GenerateAllLegalMoves()
			fmt.Printf("Legal moves (%d):\n", len(moves))
			for i, move := range moves {
				if i > 0 && i%8 == 0 {
					fmt.Println()
				}
				fmt.Printf("%-6s ", move.String())
			}
			fmt.Println()

		case "help", "h":
			fmt.Println("Commands:")
			fmt.Println("  <move>    - Make a move (e.g., e2e4, O-O)")
			fmt.Println("  ai        - Let the AI make a move")
			fmt.Println("  eval      - Show detailed position evaluation")
			fmt.Println("  depth <n> - Set AI search depth (1-10)")
			fmt.Println("  moves     - Show all legal moves")
			fmt.Println("  quit      - Exit the game")
			fmt.Println("  help      - Show this help")

		default:
			// Try to parse as a move
			move, valid := game.ParseMove(input)
			if !valid {
				fmt.Println("Invalid move format. Try e2e4 or O-O")
				continue
			}

			if game.MakeMove(move) {
				fmt.Printf("Played: %s\n", move.String())
			} else {
				fmt.Println("Illegal move!")
			}
		}
		fmt.Println()
	}
}
