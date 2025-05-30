package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Chess Engine v1.5")
	fmt.Println("===================")
	fmt.Println("Commands: move, eval, quit, moves, help")

	game := NewGame()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display current positon
		game.Board.Display()
		eval := game.EvaluatePosition()
		fmt.Printf("\nPosition evaluation %+d centipawns", eval)

		if eval > 0 {
			fmt.Println("White is better")
		} else if eval < 0 {
			fmt.Println("Black is better")
		} else {
			fmt.Println("Equal position")
		}
		fmt.Println()

		fmt.Printf("%s to move", game.GetCurrentPlayerString())

		// Check status
		if gameOver, result := game.IsGameOver(); gameOver {
			fmt.Printf("\nGame is over: %s\n", result)
			break
		}

		if game.Board.IsInCheck(game.CurrentPlayer) {
			fmt.Print("( in check)")
		}

		fmt.Print(": ")

		// Get user input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "quit", "q":
			fmt.Println("Thanks for playing!")
			return
		case "eval", "e":
			eval := game.EvaluatePosition()
			phase := game.GetGamePhase()
			isEndgame := game.isEndgame()

			fmt.Printf("Detailed evaluation\n")
			fmt.Printf("	Total Score %+d centipawns\n", eval)
			fmt.Printf("	Game Phase: %.2f (1.0=opening 0.0=endgame)\n", phase)
			fmt.Printf("	Is endgame: %t\n", isEndgame)

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
			fmt.Printf(" White material: %d\n", whiteMaterial)
			fmt.Printf(" Black material: %d\n", blackMaterial)
			fmt.Printf(" Material difference: %+d\n", whiteMaterial-blackMaterial)

		case "moves", "m":
			moves := game.GenerateAllLegalMoves()
			fmt.Printf("Legal moves (%d:)\n", len(moves))
			for i, move := range moves {
				if i > 0 && i%8 == 0 {
					fmt.Println()
				}
				fmt.Printf("%-6s", move.String())
			}
			fmt.Println()
		case "help", "h":
			fmt.Println("Commands:")
			fmt.Println(" <move> - make a move (e.g. e2e4, O-O)")
			fmt.Println(" moves - show all legal moves")
			fmt.Println(" quit - exit game")
			fmt.Println(" help - Show help")
		default:
			// Try to parse as move
			move, valid := game.ParseMove(input)
			if !valid {
				fmt.Println("Invalid format. Try e2e4 or O-O")
				continue
			}
			if game.MakeMove(move) {
				fmt.Printf("Played %s\n", move.String())
			} else {
				fmt.Println("Illegal move")
			}

		}

		fmt.Println()

	}

}
