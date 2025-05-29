package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Chess Engine v1.4")
	fmt.Println("===================")
	fmt.Println("Commands: move (e.g. e2e4), quit, moves, help")

	game := NewGame()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display current positon
		game.Board.Display()
		fmt.Printf("\n%s to move", game.GetCurrentPlayerString())

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
